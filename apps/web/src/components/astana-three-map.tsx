"use client";

import { useEffect, useRef, useState } from "react";
import { Compass, Minus, Plus } from "lucide-react";
import * as THREE from "three";
import { OrbitControls } from "three/addons/controls/OrbitControls.js";

import districtBoundaries from "@/data/astana-districts.json";
import type { DistrictResult } from "@/lib/types";

type Position = [number, number];
type Boundary = {
  properties: { id: string; name: string; osmRelationId: number; label: Position };
  geometry:
    | { type: "Polygon"; coordinates: Position[][] }
    | { type: "MultiPolygon"; coordinates: Position[][][] };
};

const BOUNDARIES = districtBoundaries.features as unknown as Boundary[];
const ORIGIN = { lon: 71.48, lat: 51.17 };
const UNITS_PER_KM = 0.25;
const MERIDIAN_KM = 111.32;
const PARALLEL_KM = MERIDIAN_KM * Math.cos(ORIGIN.lat * Math.PI / 180);
const TILE_ZOOM = 10;

type MapProps = {
  districts: DistrictResult[];
  selectedDistrictId: string | null;
  onSelectDistrict: (id: string) => void;
};

type MapRuntime = {
  camera: THREE.OrthographicCamera;
  controls: OrbitControls;
  groups: Map<string, THREE.Group>;
  meshes: THREE.Mesh[];
};

function project([lon, lat]: Position): Position {
  return [(lon - ORIGIN.lon) * PARALLEL_KM * UNITS_PER_KM, -(lat - ORIGIN.lat) * MERIDIAN_KM * UNITS_PER_KM];
}

function unproject(x: number, z: number): Position {
  return [ORIGIN.lon + x / (PARALLEL_KM * UNITS_PER_KM), ORIGIN.lat - z / (MERIDIAN_KM * UNITS_PER_KM)];
}

function ringsFor(boundary: Boundary): Position[][][] {
  return boundary.geometry.type === "Polygon"
    ? [boundary.geometry.coordinates]
    : boundary.geometry.coordinates;
}

function pathFromRing(ring: Position[], path: THREE.Path) {
  ring.forEach((coordinate, index) => {
    const [x, z] = project(coordinate);
    if (index === 0) path.moveTo(x, -z);
    else path.lineTo(x, -z);
  });
  path.closePath();
}

function makeRegion(boundary: Boundary) {
  const group = new THREE.Group();
  group.userData.districtId = boundary.properties.id;
  const meshes: THREE.Mesh[] = [];
  for (const polygon of ringsFor(boundary)) {
    const shape = new THREE.Shape();
    pathFromRing(polygon[0], shape);
    polygon.slice(1).forEach((hole) => {
      const path = new THREE.Path();
      pathFromRing(hole, path);
      shape.holes.push(path);
    });
    const geometry = new THREE.ExtrudeGeometry(shape, { depth: 0.16, bevelEnabled: true, bevelThickness: 0.018, bevelSize: -0.025, bevelSegments: 1, steps: 1 });
    geometry.rotateX(-Math.PI / 2);
    const top = new THREE.MeshStandardMaterial({ color: "#f7f3ec", transparent: true, opacity: 0.64, roughness: 0.92, depthWrite: false });
    const side = new THREE.MeshStandardMaterial({ color: "#c9c4bc", roughness: 0.92 });
    const mesh = new THREE.Mesh(geometry, [top, side]);
    mesh.userData.districtId = boundary.properties.id;
    mesh.castShadow = true;
    mesh.receiveShadow = true;
    group.add(mesh);
    meshes.push(mesh);

    const borderPoints = polygon[0].map((coordinate) => {
      const [x, z] = project(coordinate);
      return new THREE.Vector3(x, 0.22, z);
    });
    const border = new THREE.LineLoop(new THREE.BufferGeometry().setFromPoints(borderPoints), new THREE.LineBasicMaterial({ color: "#faf7f1", transparent: true, opacity: 0.9 }));
    border.userData.districtId = boundary.properties.id;
    group.add(border);
  }
  return { group, meshes };
}

function tintRegion(mesh: THREE.Mesh, selected: boolean, score: number) {
  const material = (mesh.material as THREE.Material[])[0] as THREE.MeshStandardMaterial;
  material.color.set(selected ? "#f36458" : score < 50 ? "#f4d7cb" : "#f7f3ec");
  material.opacity = selected ? 0.81 : 0.64;
}

function tileX(lon: number, zoom: number) {
  return Math.floor((lon + 180) / 360 * 2 ** zoom);
}

function tileY(lat: number, zoom: number) {
  const radians = lat * Math.PI / 180;
  return Math.floor((1 - Math.asinh(Math.tan(radians)) / Math.PI) / 2 * 2 ** zoom);
}

function tileLon(x: number, zoom: number) {
  return x / 2 ** zoom * 360 - 180;
}

function tileLat(y: number, zoom: number) {
  return Math.atan(Math.sinh(Math.PI * (1 - 2 * y / 2 ** zoom))) * 180 / Math.PI;
}

function createMapTiles(scene: THREE.Scene, renderer: THREE.WebGLRenderer, camera: THREE.OrthographicCamera, isDisposed: () => boolean) {
  const loader = new THREE.TextureLoader();
  loader.setCrossOrigin("anonymous");
  const tiles = new Map<string, THREE.Mesh | null>();
  const raycaster = new THREE.Raycaster();
  const ground = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0);
  let previousRange = "";

  return function updateTiles() {
    camera.updateMatrixWorld();
    const corners: Position[] = [];
    for (const [x, y] of [[-1, -1], [-1, 1], [1, -1], [1, 1]]) {
      raycaster.setFromCamera(new THREE.Vector2(x, y), camera);
      const point = raycaster.ray.intersectPlane(ground, new THREE.Vector3());
      if (point) corners.push(unproject(point.x, point.z));
    }
    if (corners.length !== 4) return;
    const lon = corners.map((point) => point[0]);
    const lat = corners.map((point) => point[1]);
    const minX = tileX(Math.min(...lon), TILE_ZOOM) - 1;
    const maxX = tileX(Math.max(...lon), TILE_ZOOM) + 1;
    const minY = tileY(Math.max(...lat), TILE_ZOOM) - 1;
    const maxY = tileY(Math.min(...lat), TILE_ZOOM) + 1;
    const range = `${minX}:${maxX}:${minY}:${maxY}`;
    if (range === previousRange) return;
    previousRange = range;

    for (const [key, tile] of tiles) {
      const [x, y] = key.split(":").map(Number);
      if (x >= minX && x <= maxX && y >= minY && y <= maxY) continue;
      if (tile) {
        scene.remove(tile);
        tile.geometry.dispose();
        const material = tile.material as THREE.MeshBasicMaterial;
        material.map?.dispose();
        material.dispose();
      }
      tiles.delete(key);
    }

    for (let x = minX; x <= maxX; x++) {
      for (let y = minY; y <= maxY; y++) {
      const key = `${x}:${y}`;
      if (tiles.has(key)) continue;
      tiles.set(key, null);
      const west = project([tileLon(x, TILE_ZOOM), tileLat(y + 1, TILE_ZOOM)]);
      const east = project([tileLon(x + 1, TILE_ZOOM), tileLat(y, TILE_ZOOM)]);
      const width = east[0] - west[0];
      const height = west[1] - east[1];
      const url = `https://tile.openstreetmap.org/${TILE_ZOOM}/${x}/${y}.png`;
      loader.load(url, (texture) => {
        if (isDisposed() || !tiles.has(key)) { texture.dispose(); return; }
        texture.colorSpace = THREE.SRGBColorSpace;
        texture.anisotropy = Math.min(8, renderer.capabilities.getMaxAnisotropy());
        const material = new THREE.MeshBasicMaterial({ map: texture, transparent: true, opacity: 0.72, depthWrite: false });
        const tile = new THREE.Mesh(new THREE.PlaneGeometry(width, height), material);
        tile.rotation.x = -Math.PI / 2;
        tile.position.set((west[0] + east[0]) / 2, -0.07, (west[1] + east[1]) / 2);
        tiles.set(key, tile);
        scene.add(tile);
      }, undefined, () => { tiles.delete(key); });
      }
    }
  };
}

export function AstanaThreeMap({ districts, selectedDistrictId, onSelectDistrict }: MapProps) {
  const mountRef = useRef<HTMLDivElement>(null);
  const labelsRef = useRef<Map<string, HTMLButtonElement>>(new Map());
  const runtimeRef = useRef<MapRuntime | null>(null);
  const selectionRef = useRef(onSelectDistrict);
  const [available, setAvailable] = useState(true);

  useEffect(() => { selectionRef.current = onSelectDistrict; }, [onSelectDistrict]);

  useEffect(() => {
    const mount = mountRef.current;
    if (!mount) return;
    let renderer: THREE.WebGLRenderer;
    try {
      renderer = new THREE.WebGLRenderer({ antialias: true, alpha: false, powerPreference: "high-performance" });
    } catch {
      window.setTimeout(() => setAvailable(false), 0);
      return;
    }
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    renderer.setClearColor("#eee9e3");
    renderer.shadowMap.enabled = true;
    renderer.shadowMap.type = THREE.PCFSoftShadowMap;
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    mount.appendChild(renderer.domElement);

    const scene = new THREE.Scene();
    scene.background = new THREE.Color("#eee9e3");
    const camera = new THREE.OrthographicCamera(-10, 10, 8, -8, 0.1, 100);
    const target = new THREE.Vector3(0, 0, 1);
    camera.position.set(12, 17, 19);
    camera.lookAt(target);
    camera.updateProjectionMatrix();

    const controls = new OrbitControls(camera, renderer.domElement);
    controls.enableDamping = !window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    controls.dampingFactor = 0.09;
    controls.minPolarAngle = 0.42;
    controls.maxPolarAngle = 1.18;
    controls.minZoom = 0.85;
    controls.maxZoom = 2.4;
    controls.enablePan = true;
    controls.screenSpacePanning = true;
    controls.target.copy(target);

    scene.add(new THREE.HemisphereLight("#ffffff", "#c8bdb7", 2.5));
    const sun = new THREE.DirectionalLight("#fff8ee", 2.1);
    sun.position.set(-6, 13, 8);
    sun.castShadow = true;
    sun.shadow.mapSize.set(2048, 2048);
    sun.shadow.camera.left = -15;
    sun.shadow.camera.right = 15;
    sun.shadow.camera.top = 15;
    sun.shadow.camera.bottom = -15;
    scene.add(sun);

    const ground = new THREE.Mesh(new THREE.PlaneGeometry(200, 200), new THREE.MeshStandardMaterial({ color: "#eee9e3", roughness: 1 }));
    ground.rotation.x = -Math.PI / 2;
    ground.position.y = -0.18;
    ground.receiveShadow = true;
    scene.add(ground);

    let disposed = false;
    const updateTiles = createMapTiles(scene, renderer, camera, () => disposed);
    controls.addEventListener("change", updateTiles);

    const groups = new Map<string, THREE.Group>();
    const meshes: THREE.Mesh[] = [];
    for (const boundary of BOUNDARIES) {
      const region = makeRegion(boundary);
      groups.set(boundary.properties.id, region.group);
      meshes.push(...region.meshes);
      scene.add(region.group);
    }

    const raycaster = new THREE.Raycaster();
    const pointer = new THREE.Vector2();
    const canvas = renderer.domElement;
    let downAt: [number, number] | null = null;
    let moved = false;
    let frame = 0;

    function hit(clientX: number, clientY: number) {
      const rect = canvas.getBoundingClientRect();
      pointer.set(((clientX - rect.left) / rect.width) * 2 - 1, -((clientY - rect.top) / rect.height) * 2 + 1);
      raycaster.setFromCamera(pointer, camera);
      return raycaster.intersectObjects(meshes, false)[0]?.object.userData.districtId as string | undefined;
    }
    function onPointerDown(event: PointerEvent) {
      downAt = [event.clientX, event.clientY];
      moved = false;
    }
    function onPointerMove(event: PointerEvent) {
      if (downAt && Math.hypot(event.clientX - downAt[0], event.clientY - downAt[1]) > 5) moved = true;
      canvas.style.cursor = hit(event.clientX, event.clientY) ? "pointer" : "grab";
    }
    function onPointerUp(event: PointerEvent) {
      if (downAt && !moved) {
        const district = hit(event.clientX, event.clientY);
        if (district) selectionRef.current(district);
      }
      downAt = null;
    }
    canvas.addEventListener("pointerdown", onPointerDown);
    canvas.addEventListener("pointermove", onPointerMove);
    canvas.addEventListener("pointerup", onPointerUp);

    function resize() {
      if (!mount) return;
      const { width, height } = mount.getBoundingClientRect();
      if (!width || !height) return;
      const halfHeight = 8.4;
      const aspect = width / height;
      camera.left = -halfHeight * aspect;
      camera.right = halfHeight * aspect;
      camera.top = halfHeight;
      camera.bottom = -halfHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height, false);
      updateTiles();
    }
    const observer = new ResizeObserver(resize);
    observer.observe(mount);
    resize();

    function render() {
      controls.update();
      for (const boundary of BOUNDARIES) {
        const label = labelsRef.current.get(boundary.properties.id);
        const group = groups.get(boundary.properties.id);
        if (!label || !group || !mount) continue;
        const [x, z] = project(boundary.properties.label);
        const projected = new THREE.Vector3(x, group.position.y + 0.64, z).project(camera);
        label.style.left = `${(projected.x * 0.5 + 0.5) * mount.clientWidth}px`;
        label.style.top = `${(-projected.y * 0.5 + 0.5) * mount.clientHeight}px`;
        label.style.opacity = projected.z < 1 ? "1" : "0";
      }
      renderer.render(scene, camera);
      frame = window.requestAnimationFrame(render);
    }
    runtimeRef.current = { camera, controls, groups, meshes };
    render();

    return () => {
      disposed = true;
      window.cancelAnimationFrame(frame);
      observer.disconnect();
      controls.removeEventListener("change", updateTiles);
      canvas.removeEventListener("pointerdown", onPointerDown);
      canvas.removeEventListener("pointermove", onPointerMove);
      canvas.removeEventListener("pointerup", onPointerUp);
      controls.dispose();
      scene.traverse((object) => {
        if (object instanceof THREE.Mesh) {
          object.geometry.dispose();
          const materials = Array.isArray(object.material) ? object.material : [object.material];
          materials.forEach((material) => {
            if (material instanceof THREE.MeshBasicMaterial && material.map) material.map.dispose();
            material.dispose();
          });
        } else if (object instanceof THREE.Line) {
          object.geometry.dispose();
          (object.material as THREE.Material).dispose();
        }
      });
      renderer.dispose();
      renderer.domElement.remove();
      runtimeRef.current = null;
    };
  }, []);

  useEffect(() => {
    const runtime = runtimeRef.current;
    if (!runtime) return;
    for (const boundary of BOUNDARIES) {
      const group = runtime.groups.get(boundary.properties.id);
      if (group) group.position.y = selectedDistrictId === boundary.properties.id ? 0.2 : 0;
    }
    for (const mesh of runtime.meshes) {
      const id = mesh.userData.districtId as string;
      const score = districts.find((district) => district.id === id)?.scoreAfter ?? 55;
      tintRegion(mesh, selectedDistrictId === id, score);
    }
  }, [districts, selectedDistrictId]);

  function zoom(factor: number) {
    const runtime = runtimeRef.current;
    if (!runtime) return;
    runtime.camera.zoom = THREE.MathUtils.clamp(runtime.camera.zoom * factor, 0.85, 2.4);
    runtime.camera.updateProjectionMatrix();
  }

  function resetView() {
    const runtime = runtimeRef.current;
    if (!runtime) return;
    runtime.camera.position.set(12, 17, 19);
    runtime.camera.zoom = 1;
    runtime.camera.updateProjectionMatrix();
    runtime.controls.target.set(0, 0, 1);
    runtime.controls.update();
  }

  return <div className="astana-map" aria-label="Интерактивная карта шести районов Астаны">
    <div className="astana-map-canvas" ref={mountRef} aria-hidden="true" />
    {available ? BOUNDARIES.map((boundary) => {
      const district = districts.find((item) => item.id === boundary.properties.id);
      return <button
        className={`astana-map-label ${selectedDistrictId === boundary.properties.id ? "is-selected" : ""}`}
        key={boundary.properties.id}
        ref={(node) => { if (node) labelsRef.current.set(boundary.properties.id, node); else labelsRef.current.delete(boundary.properties.id); }}
        type="button"
        onClick={() => onSelectDistrict(boundary.properties.id)}
        aria-pressed={selectedDistrictId === boundary.properties.id}
      ><strong>{district?.name ?? boundary.properties.name}</strong><span>{district ? district.scoreAfter.toFixed(1) : "—"}</span></button>;
    }) : <div className="astana-map-fallback">{BOUNDARIES.map((boundary) => <button key={boundary.properties.id} type="button" onClick={() => onSelectDistrict(boundary.properties.id)}>{boundary.properties.name}</button>)}</div>}
    <div className="astana-map-tools" aria-label="Управление картой">
      <button type="button" onClick={() => zoom(1.18)} aria-label="Приблизить карту"><Plus size={19} /></button>
      <button type="button" onClick={() => zoom(1 / 1.18)} aria-label="Отдалить карту"><Minus size={19} /></button>
      <button type="button" onClick={resetView} aria-label="Сбросить вид карты"><Compass size={18} /></button>
    </div>
    <a className="astana-map-attribution" href="https://www.openstreetmap.org/copyright" target="_blank" rel="noreferrer">© OpenStreetMap contributors</a>
  </div>;
}

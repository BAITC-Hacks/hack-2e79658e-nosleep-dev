import type { DistrictResult } from "@/lib/types";

const shapes = [
  { id: "baikonur", points: "88,48 210,30 250,104 224,178 112,160 62,102", label: [148, 99] },
  { id: "saryarka", points: "62,106 112,164 222,182 198,258 92,270 28,202", label: [115, 216] },
  { id: "nura", points: "214,30 356,54 390,144 322,196 230,176 254,103", label: [310, 112] },
  { id: "yesil", points: "226,182 320,202 344,294 250,342 166,286 202,260", label: [258, 264] },
  { id: "almaty", points: "90,276 164,290 246,348 208,420 72,390 30,326", label: [136, 348] },
];

function fillFor(score: number) {
  if (score < 50) return "#f36458";
  if (score < 56) return "#e7a465";
  if (score < 61) return "#d7d9b0";
  return "#74d8a3";
}

export function CityMap({ districts }: { districts: DistrictResult[] }) {
  return (
    <div className="map-shell" aria-label="Стилизованная карта районов Астаны">
      <svg viewBox="0 0 420 450" role="img" aria-labelledby="map-title map-description">
        <title id="map-title">Карта качества жизни по районам</title>
        <desc id="map-description">Пять районов окрашены в зависимости от текущего балла.</desc>
        {shapes.map((shape) => {
          const district = districts.find((item) => item.id === shape.id);
          const score = district?.scoreAfter ?? 0;
          return (
            <g key={shape.id} className="district-shape">
              <polygon points={shape.points} fill={fillFor(score)} />
              <text x={shape.label[0]} y={shape.label[1]} textAnchor="middle">{district?.name ?? shape.id}</text>
              <text className="district-score" x={shape.label[0]} y={shape.label[1] + 22} textAnchor="middle">{score.toFixed(1)}</text>
            </g>
          );
        })}
      </svg>
      <div className="map-legend" aria-hidden="true">
        <span><i className="critical" />критично</span>
        <span><i className="watch" />внимание</span>
        <span><i className="stable" />стабильно</span>
      </div>
    </div>
  );
}

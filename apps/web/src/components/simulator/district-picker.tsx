"use client";

import { useEffect } from "react";
import { Building2, MapPin, X } from "lucide-react";

import type { District, Initiative } from "@/lib/types";

type DistrictPickerProps = {
  districts: District[];
  initiative: Initiative;
  onCancel: () => void;
  onSelect: (districtId: string) => void;
};

export function DistrictPicker({ districts, initiative, onCancel, onSelect }: DistrictPickerProps) {
  useEffect(() => {
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") onCancel();
    };
    window.addEventListener("keydown", closeOnEscape);
    return () => window.removeEventListener("keydown", closeOnEscape);
  }, [onCancel]);

  return (
    <div className="district-picker-backdrop" role="presentation" onMouseDown={(event) => { if (event.currentTarget === event.target) onCancel(); }}>
      <section className="district-picker" role="dialog" aria-modal="true" aria-labelledby="district-picker-title">
        <div className="district-picker-heading">
          <div><span className="cockpit-kicker">ВЫБОР РАЙОНА</span><h2 id="district-picker-title">Где применить меру?</h2></div>
          <button type="button" onClick={onCancel} aria-label="Закрыть выбор района" autoFocus><X size={20} /></button>
        </div>
        <p>{initiative.name}</p>
        <div className="district-options">
          {districts.map((district) => (
            <button type="button" key={district.id} onClick={() => onSelect(district.id)}>
              <span><MapPin size={18} /><strong>{district.name}</strong></span>
              <small>{district.profile}</small>
              <Building2 size={17} />
            </button>
          ))}
        </div>
      </section>
    </div>
  );
}

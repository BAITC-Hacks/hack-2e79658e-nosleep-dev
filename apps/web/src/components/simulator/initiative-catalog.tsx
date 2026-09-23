import { Check, MapPin, Plus, Search } from "lucide-react";

import type { Direction, Initiative } from "@/lib/types";

export const directionLabels: Record<Direction, string> = {
  transport: "Транспорт",
  ecology: "Экология",
  social: "Соцсфера",
  safety: "Безопасность",
  services: "Сервисы",
};

type DirectionFilter = "all" | Direction;

type InitiativeCatalogProps = {
  activeDirection: DirectionFilter;
  budgetRemaining: number;
  directionCounts: Partial<Record<Direction, number>>;
  initiatives: Initiative[];
  query: string;
  requiredDecisions: number;
  selectedIds: Set<string>;
  selectedCount: number;
  onDirectionChange: (direction: DirectionFilter) => void;
  onQueryChange: (query: string) => void;
  onRequestAdd: (initiative: Initiative) => void;
  onRemove: (initiativeId: string) => void;
};

const filters: DirectionFilter[] = ["all", "transport", "ecology", "social", "safety", "services"];

export function InitiativeCatalog(props: InitiativeCatalogProps) {
  const filtered = props.initiatives.filter((initiative) => {
    const matchesDirection = props.activeDirection === "all" || initiative.direction === props.activeDirection;
    const matchesQuery = initiative.name.toLocaleLowerCase("ru").includes(props.query.trim().toLocaleLowerCase("ru"));
    return matchesDirection && matchesQuery;
  });

  return (
    <section className="cockpit-panel catalog-workbench" aria-labelledby="catalog-title">
      <div className="cockpit-panel-heading">
        <div><span className="cockpit-kicker">ШАГ 01</span><h2 id="catalog-title">Выберите инициативы</h2></div>
        <span className="catalog-count">{filtered.length} мер</span>
      </div>

      <label className="catalog-search">
        <Search size={17} />
        <span className="sr-only">Поиск инициатив</span>
        <input value={props.query} onChange={(event) => props.onQueryChange(event.target.value)} placeholder="Найти инициативу" />
      </label>

      <div className="catalog-filters" aria-label="Фильтр по направлениям">
        {filters.map((filter) => (
          <button className={props.activeDirection === filter ? "is-active" : ""} type="button" key={filter} onClick={() => props.onDirectionChange(filter)}>
            {filter === "all" ? "Все" : directionLabels[filter]}
          </button>
        ))}
      </div>

      <div className="cockpit-initiative-list">
        {filtered.length === 0 && <div className="catalog-empty">По вашему запросу ничего не найдено.</div>}
        {filtered.map((initiative) => {
          const selected = props.selectedIds.has(initiative.id);
          const directionFull = (props.directionCounts[initiative.direction] ?? 0) >= 2;
          const noSlots = props.selectedCount >= props.requiredDecisions;
          const tooExpensive = initiative.cost > props.budgetRemaining;
          const disabled = !selected && (directionFull || noSlots || tooExpensive);
          const reason = selected ? undefined : noSlots ? "Все пять слотов заняты" : directionFull ? "Лимит: две меры направления" : tooExpensive ? "Недостаточно бюджета" : undefined;

          return (
            <article className={`cockpit-initiative ${selected ? "is-selected" : ""} ${disabled ? "is-disabled" : ""}`} key={initiative.id}>
              <div className="cockpit-initiative-top">
                <span className="initiative-code">{initiative.id}</span>
                <span className="initiative-direction">{directionLabels[initiative.direction]}</span>
                <strong>{initiative.cost}<small> ед.</small></strong>
              </div>
              <h3>{initiative.name}</h3>
              <div className="cockpit-initiative-bottom">
                <span>{initiative.type === "district" ? <><MapPin size={13} /> Для района</> : "Для всего города"}</span>
                {selected ? (
                  <button type="button" onClick={() => props.onRemove(initiative.id)} aria-label={`Убрать ${initiative.name}`}><Check size={15} /> Выбрано</button>
                ) : (
                  <button type="button" onClick={() => props.onRequestAdd(initiative)} disabled={disabled} title={reason} aria-label={`Добавить ${initiative.name}`}><Plus size={15} /> Добавить</button>
                )}
              </div>
              {reason && <p className="initiative-disabled-reason">{reason}</p>}
            </article>
          );
        })}
      </div>
    </section>
  );
}

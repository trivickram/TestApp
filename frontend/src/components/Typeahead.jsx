import { useState, useRef } from "react";

export function Typeahead({
  search,
  labelFn,
  placeholder,
  onSelect,
  disabled,
}) {
  const [q, setQ] = useState("");
  const [hits, setHits] = useState([]);
  const timer = useRef(null);

  const onChange = (e) => {
    const v = e.target.value;
    setQ(v);
    clearTimeout(timer.current);
    if (v.length < 2) {
      setHits([]);
      return;
    }
    timer.current = setTimeout(async () => {
      const d = await search(v);
      setHits(d ?? []);
    }, 250);
  };

  const pick = (item) => {
    setQ("");
    setHits([]);
    onSelect(item);
  };

  return (
    <div className="ta">
      <input
        value={q}
        onChange={onChange}
        placeholder={placeholder}
        disabled={disabled}
      />
      {hits.length > 0 && (
        <ul className="ta-list">
          {hits.map((h) => (
            <li key={h.id} onMouseDown={() => pick(h)}>
              {labelFn(h)}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

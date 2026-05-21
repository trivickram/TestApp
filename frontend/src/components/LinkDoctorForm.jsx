import { useState } from "react";
import { api } from "../api";
import { Typeahead } from "./Typeahead";

export function LinkDoctorForm({ clinicId, onLinked }) {
  const [selected, setSelected] = useState(null);
  const [error, setError] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    if (!selected) return;
    const r = await api.linkDoctor(clinicId, selected.id);
    if (!r.ok) {
      setError((await r.json()).error ?? "error");
      return;
    }
    setSelected(null);
    onLinked();
  };

  return (
    <fieldset disabled={!clinicId}>
      <legend>Link Existing Doctor</legend>
      {error && (
        <p className="error">
          {error}{" "}
          <button type="button" onClick={() => setError("")}>
            ×
          </button>
        </p>
      )}
      <form className="row" onSubmit={submit}>
        {selected ? (
          <span className="chip">
            {selected.name} — {selected.specialization}
            <button type="button" onClick={() => setSelected(null)}>
              ×
            </button>
          </span>
        ) : (
          <Typeahead
            search={api.searchDoctors}
            labelFn={(d) => `${d.name} — ${d.specialization}`}
            placeholder="search by name..."
            onSelect={setSelected}
            disabled={!clinicId}
          />
        )}
        <button type="submit" disabled={!selected}>
          link to clinic
        </button>
      </form>
    </fieldset>
  );
}

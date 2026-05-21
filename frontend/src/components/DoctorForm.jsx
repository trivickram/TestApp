import { useState } from "react";
import { api } from "../api";

export function DoctorForm({ clinicId, onAdded }) {
  const [name, setName] = useState("");
  const [spec, setSpec] = useState("");
  const [error, setError] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    const r = await api.createDoctor({ name, specialization: spec });
    if (!r.ok) {
      setError((await r.json()).error ?? "error");
      return;
    }
    const doc = await r.json();
    const r2 = await api.linkDoctor(clinicId, doc.id);
    if (!r2.ok) {
      setError((await r2.json()).error ?? "error");
      return;
    }
    setName("");
    setSpec("");
    onAdded();
  };

  return (
    <fieldset disabled={!clinicId}>
      <legend>Create Doctor</legend>
      {error && (
        <p className="error">
          {error}{" "}
          <button type="button" onClick={() => setError("")}>
            ×
          </button>
        </p>
      )}
      <form className="row" onSubmit={submit}>
        <input
          placeholder="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
        <input
          placeholder="specialization"
          value={spec}
          onChange={(e) => setSpec(e.target.value)}
          required
        />
        <button type="submit">create &amp; link</button>
      </form>
    </fieldset>
  );
}

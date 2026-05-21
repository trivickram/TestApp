import { useState } from "react";
import { api } from "../api";

export function PatientForm({ clinicId, onAdded }) {
  const [name, setName] = useState("");
  const [age, setAge] = useState("");
  const [error, setError] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    const r = await api.createPatient({ name, age: +age });
    if (!r.ok) {
      setError((await r.json()).error ?? "error");
      return;
    }
    setName("");
    setAge("");
    onAdded();
  };

  return (
    <fieldset disabled={!clinicId}>
      <legend>New Patient</legend>
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
          type="number"
          placeholder="age"
          value={age}
          onChange={(e) => setAge(e.target.value)}
          required
          className="short"
        />
        <button type="submit">add</button>
      </form>
    </fieldset>
  );
}

import { useState } from "react";
import { api } from "../api";
import { Typeahead } from "./Typeahead";

export function ScheduleForm({ clinicId, doctors, patients, onScheduled }) {
  const [doctor, setDoctor] = useState(null);
  const [patient, setPatient] = useState(null);
  const [time, setTime] = useState("");
  const [error, setError] = useState("");

  const searchDoctors = (q) =>
    Promise.resolve(
      (doctors ?? []).filter((d) =>
        d.name.toLowerCase().includes(q.toLowerCase()),
      ),
    );

  const searchPatients = (q) =>
    Promise.resolve(
      (patients ?? []).filter((p) =>
        p.name.toLowerCase().includes(q.toLowerCase()),
      ),
    );

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    if (!doctor || !patient) {
      setError("select a doctor and patient");
      return;
    }
    const r = await api.scheduleAppointment({
      clinic_id: +clinicId,
      doctor_id: doctor.id,
      patient_id: patient.id,
      scheduled_at: time.replace("T", " ") + ":00",
    });
    if (!r.ok) {
      setError((await r.json()).error ?? "error");
      return;
    }
    setDoctor(null);
    setPatient(null);
    setTime("");
    onScheduled();
  };

  return (
    <fieldset disabled={!clinicId}>
      <legend>Schedule Appointment</legend>
      {error && (
        <p className="error">
          {error}{" "}
          <button type="button" onClick={() => setError("")}>
            ×
          </button>
        </p>
      )}
      <form className="row" onSubmit={submit}>
        {doctor ? (
          <span className="chip">
            {doctor.name}
            <button type="button" onClick={() => setDoctor(null)}>
              ×
            </button>
          </span>
        ) : (
          <Typeahead
            search={searchDoctors}
            labelFn={(d) => `${d.name} — ${d.specialization}`}
            placeholder="search doctor..."
            onSelect={setDoctor}
            disabled={!clinicId}
          />
        )}
        {patient ? (
          <span className="chip">
            {patient.name}
            <button type="button" onClick={() => setPatient(null)}>
              ×
            </button>
          </span>
        ) : (
          <Typeahead
            search={searchPatients}
            labelFn={(p) => `${p.name} (age ${p.age})`}
            placeholder="search patient..."
            onSelect={setPatient}
            disabled={!clinicId}
          />
        )}
        <input
          type="datetime-local"
          value={time}
          onChange={(e) => setTime(e.target.value)}
          required
        />
        <button type="submit">schedule</button>
      </form>
    </fieldset>
  );
}

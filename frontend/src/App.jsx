import { useState } from "react";

const apiBase = "http://localhost:8080";

export default function App() {
  const [patientName, setPatientName] = useState("");
  const [patientAge, setPatientAge] = useState("");
  const [patient, setPatient] = useState(null);

  const [appointmentDoctor, setAppointmentDoctor] = useState("");
  const [appointmentTime, setAppointmentTime] = useState("");
  const [appointment, setAppointment] = useState(null);

  const [appointmentID, setAppointmentID] = useState("");
  const [status, setStatus] = useState(null);
  const [updateStatusValue, setUpdateStatusValue] = useState("CONFIRMED");
  const [error, setError] = useState("");

  async function createPatient(e) {
    e.preventDefault();
    setError("");
    setStatus(null);

    const response = await fetch(`${apiBase}/patients`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: patientName,
        age: Number(patientAge),
      }),
    });

    const data = await response.json();
    if (!response.ok) {
      setError(data.error || "failed to create patient");
      return;
    }

    setPatient(data);
  }

  async function scheduleAppointment(e) {
    e.preventDefault();
    setError("");
    setStatus(null);

    const patientID = patient?.id;
    if (!patientID) {
      setError("create a patient first");
      return;
    }

    const response = await fetch(`${apiBase}/appointments`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        patient_id: patientID,
        doctor: appointmentDoctor,
        scheduled_at: appointmentTime,
      }),
    });

    const data = await response.json();
    if (!response.ok) {
      setError(data.error || "failed to schedule appointment");
      return;
    }

    setAppointment(data);
    setAppointmentID(data.id);
  }

  async function getStatus(e) {
    e.preventDefault();
    setError("");

    const response = await fetch(
      `${apiBase}/appointments/status?appointment_id=${encodeURIComponent(appointmentID)}`,
    );
    const data = await response.json();
    if (!response.ok) {
      setError(data.error || "failed to get status");
      return;
    }

    setStatus(data);
  }

  async function updateStatus(e) {
    e.preventDefault();
    setError("");

    const response = await fetch(`${apiBase}/appointments/status/update`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        appointment_id: appointmentID,
        status: updateStatusValue,
      }),
    });

    const data = await response.json();
    if (!response.ok) {
      setError(data.error || "failed to update status");
      return;
    }

    setStatus(data);
  }

  return (
    <main className="page">
      <div className="panel">
        <h1>Hospital Management</h1>

        <form onSubmit={createPatient} className="card">
          <h2>Create Patient</h2>
          <input
            placeholder="Name"
            value={patientName}
            onChange={(e) => setPatientName(e.target.value)}
          />
          <input
            type="number"
            placeholder="Age"
            value={patientAge}
            onChange={(e) => setPatientAge(e.target.value)}
          />
          <button type="submit">Create</button>
          {patient && <pre>{JSON.stringify(patient, null, 2)}</pre>}
        </form>

        <form onSubmit={scheduleAppointment} className="card">
          <h2>Schedule Appointment</h2>
          <input
            placeholder="Doctor"
            value={appointmentDoctor}
            onChange={(e) => setAppointmentDoctor(e.target.value)}
          />
          <input
            type="datetime-local"
            value={appointmentTime}
            onChange={(e) => setAppointmentTime(e.target.value)}
          />
          <button type="submit">Schedule</button>
          {appointment && <pre>{JSON.stringify(appointment, null, 2)}</pre>}
        </form>

        <form onSubmit={getStatus} className="card">
          <h2>Get Appointment Status</h2>
          <input
            placeholder="Appointment ID"
            // value={appointmentID}
            onChange={(e) => setAppointmentID(e.target.value)}
          />
          <button type="submit">Get Status</button>
          {status && <pre>{JSON.stringify(status, null, 2)}</pre>}
        </form>

        <form onSubmit={updateStatus} className="card">
          <h2>Update Appointment Status</h2>
          <input
            placeholder="Appointment ID"
            // value={appointmentID}
            onChange={(e) => setAppointmentID(e.target.value)}
          />
          <select
            value={updateStatusValue}
            onChange={(e) => setUpdateStatusValue(e.target.value)}
          >
            <option value="CONFIRMED">CONFIRMED</option>
            <option value="IN_PROGRESS">IN_PROGRESS</option>
            <option value="COMPLETED">COMPLETED</option>
            <option value="CANCELLED">CANCELLED</option>
          </select>
          <button type="submit">Update Status</button>
          {status && <pre>{JSON.stringify(status, null, 2)}</pre>}
        </form>

        {error && <p className="error">{error}</p>}
      </div>
    </main>
  );
}

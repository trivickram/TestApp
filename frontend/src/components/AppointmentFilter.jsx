import { useState } from "react";

export function AppointmentFilter({
  clinics,
  filterDoctors,
  filterPatients,
  onClinicChange,
  onFilter,
}) {
  const [fClinic, setFClinic] = useState("");
  const [fDoctor, setFDoctor] = useState("");
  const [fPatient, setFPatient] = useState("");
  const [fDate, setFDate] = useState("");

  const handleClinicChange = (id) => {
    setFClinic(id);
    setFDoctor("");
    setFPatient("");
    onClinicChange(id);
    onFilter({ clinic_id: id });
  };

  const apply = () =>
    onFilter({
      clinic_id: fClinic,
      doctor_id: fDoctor,
      patient_id: fPatient,
      date: fDate,
    });

  const clear = () => {
    setFClinic("");
    setFDoctor("");
    setFPatient("");
    setFDate("");
    onClinicChange("");
    onFilter({});
  };

  return (
    <div className="row">
      <label>Filter</label>
      <select
        value={fClinic}
        onChange={(e) => handleClinicChange(e.target.value)}
      >
        <option value="">all clinics</option>
        {clinics.map((c) => (
          <option key={c.id} value={c.id}>
            {c.name}
          </option>
        ))}
      </select>
      <select
        value={fDoctor}
        onChange={(e) => setFDoctor(e.target.value)}
        disabled={!fClinic}
      >
        <option value="">all doctors</option>
        {filterDoctors.map((d) => (
          <option key={d.id} value={d.id}>
            {d.name}
          </option>
        ))}
      </select>
      <select
        value={fPatient}
        onChange={(e) => setFPatient(e.target.value)}
        disabled={!fClinic}
      >
        <option value="">all patients</option>
        {filterPatients.map((p) => (
          <option key={p.id} value={p.id}>
            {p.name}
          </option>
        ))}
      </select>
      <input
        type="date"
        value={fDate}
        onChange={(e) => setFDate(e.target.value)}
      />
      <button onClick={apply}>filter</button>
      <button onClick={clear}>clear</button>
    </div>
  );
}

import { useState, useEffect } from "react";
import "./styles.css";
import { api } from "./api";
import { PatientForm } from "./components/PatientForm";
import { DoctorForm } from "./components/DoctorForm";
import { LinkDoctorForm } from "./components/LinkDoctorForm";
import { ScheduleForm } from "./components/ScheduleForm";
import { AppointmentFilter } from "./components/AppointmentFilter";
import { AppointmentTable } from "./components/AppointmentTable";

export default function App() {
  const [clinics, setClinics] = useState([]);
  const [clinicId, setClinicId] = useState("");
  const [doctors, setDoctors] = useState([]);
  const [patients, setPatients] = useState([]);
  const [appointments, setAppointments] = useState([]);

  const [filterDoctors, setFilterDoctors] = useState([]);
  const [filterPatients, setFilterPatients] = useState([]);

  useEffect(() => {
    api.getClinics().then(setClinics);
  }, []);

  useEffect(() => {
    if (!clinicId) {
      setDoctors([]);
      setPatients([]);
      return;
    }
    api.getClinicDoctors(clinicId).then(setDoctors);
    api.getClinicPatients(clinicId).then(setPatients);
  }, [clinicId]);

  const refreshDoctors = () => api.getClinicDoctors(clinicId).then(setDoctors);
  const refreshPatients = () =>
    api.getClinicPatients(clinicId).then(setPatients);

  const handleFilterClinicChange = (id) => {
    if (!id) {
      setFilterDoctors([]);
      setFilterPatients([]);
      return;
    }
    api.getClinicDoctors(id).then(setFilterDoctors);
    api.getClinicPatients(id).then(setFilterPatients);
  };

  const handleFilter = (params) =>
    api.getAppointments(params).then(setAppointments);

  const handleScheduled = () => handleFilterWithCache({});

  const [lastFilterParams, setLastFilterParams] = useState({});

  const handleFilterWithCache = (params) => {
    setLastFilterParams(params);
    return api.getAppointments(params).then(setAppointments);
  };

  const handleStatusChange = async (id, newStatus) => {
    const res = await api.updateAppointmentStatus(id, newStatus);
    if (!res.ok) return;
    api.getAppointments(lastFilterParams).then(setAppointments);
  };

  // Merge all loaded doctor/patient lists for name resolution in table
  const merge = (a, b) => {
    const map = new Map(a.map((x) => [x.id, x]));
    b.forEach((x) => map.set(x.id, x));
    return [...map.values()];
  };
  const allDoctors = merge(doctors, filterDoctors);
  const allPatients = merge(patients, filterPatients);

  return (
    <div className="app">
      <h1>Hospital</h1>

      <div className="row">
        <label>Clinic</label>
        <select value={clinicId} onChange={(e) => setClinicId(e.target.value)}>
          <option value="">select...</option>
          {clinics.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
      </div>

      <hr />

      <PatientForm clinicId={clinicId} onAdded={refreshPatients} />
      <DoctorForm clinicId={clinicId} onAdded={refreshDoctors} />
      <LinkDoctorForm clinicId={clinicId} onLinked={refreshDoctors} />
      <ScheduleForm
        clinicId={clinicId}
        doctors={doctors}
        patients={patients}
        onScheduled={handleScheduled}
      />

      <hr />

      <AppointmentFilter
        clinics={clinics}
        filterDoctors={filterDoctors}
        filterPatients={filterPatients}
        onClinicChange={handleFilterClinicChange}
        onFilter={handleFilterWithCache}
      />
      <AppointmentTable
        appointments={appointments}
        clinics={clinics}
        doctors={allDoctors}
        patients={allPatients}
        onStatusChange={handleStatusChange}
      />
    </div>
  );
}

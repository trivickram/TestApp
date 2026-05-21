export function AppointmentTable({
  appointments,
  clinics,
  doctors,
  patients,
  onStatusChange,
}) {
  const find = (list, id) => list.find((x) => x.id == id);

  const clinicName = (id) => find(clinics, id)?.name ?? id;
  const doctorName = (id) => find(doctors, id)?.name ?? id;
  const patientName = (id) => find(patients, id)?.name ?? id;

  return (
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Clinic</th>
          <th>Doctor</th>
          <th>Patient</th>
          <th>Time</th>
          <th>Status</th>
          <th>Actions</th>
        </tr>
      </thead>
      <tbody>
        {appointments.length === 0 ? (
          <tr>
            <td colSpan={7} className="empty">
              no appointments
            </td>
          </tr>
        ) : (
          appointments.map((a) => (
            <tr key={a.id}>
              <td>{a.id}</td>
              <td>{clinicName(a.clinic_id)}</td>
              <td>{doctorName(a.doctor_id)}</td>
              <td>{patientName(a.patient_id)}</td>
              <td>{a.scheduled_at}</td>
              <td>{a.status}</td>
              <td>
                {a.status === "SCHEDULED" && (
                  <>
                    <button onClick={() => onStatusChange(a.id, "COMPLETED")}>
                      Complete
                    </button>{" "}
                    <button onClick={() => onStatusChange(a.id, "CANCELLED")}>
                      Cancel
                    </button>
                  </>
                )}
              </td>
            </tr>
          ))
        )}
      </tbody>
    </table>
  );
}

const BASE = "http://localhost:8080";

const post = (url, body) =>
  fetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

const json = (res) => res.json().then((d) => d ?? []);

export const api = {
  getClinics: () => fetch(`${BASE}/clinics`).then(json),

  getClinicDoctors: (id) => fetch(`${BASE}/clinics/${id}/doctors`).then(json),
  getClinicPatients: (id) => fetch(`${BASE}/clinics/${id}/patients`).then(json),

  searchDoctors: (q) =>
    fetch(`${BASE}/doctors?q=${encodeURIComponent(q)}`).then(json),
  searchPatients: (q) =>
    fetch(`${BASE}/patients?q=${encodeURIComponent(q)}`).then(json),

  getAppointments: (params) => {
    const p = new URLSearchParams();
    Object.entries(params).forEach(([k, v]) => {
      if (v) p.set(k, v);
    });
    return fetch(`${BASE}/appointments?${p}`).then(json);
  },

  createPatient: (data) => post(`${BASE}/patients`, data),
  createDoctor: (data) => post(`${BASE}/doctors`, data),
  linkDoctor: (clinicId, doctorId) =>
    post(`${BASE}/clinics/${clinicId}/doctors`, { doctor_id: doctorId }),
  scheduleAppointment: (data) => post(`${BASE}/appointments`, data),

  updateAppointmentStatus: (id, status) =>
    fetch(`${BASE}/appointments/${id}/status`, {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ status }),
    }),

  saveConsultation: (data) => post(`${BASE}/consultations`, data),

  getConsultation: (appointmentId) =>
    fetch(`${BASE}/appointments/${appointmentId}/consultation`),
};

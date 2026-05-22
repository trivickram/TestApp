import { useState, useEffect } from "react";
import { api } from "../api";

export function ConsultationPad({
  appointment,
  doctorName,
  patientName,
  onClose,
}) {
  const readOnly = appointment.status === "COMPLETED";

  const [symptoms, setSymptoms] = useState("");
  const [prescription, setPrescription] = useState("");
  const [notes, setNotes] = useState("");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    api.getConsultation(appointment.id).then(async (res) => {
      if (res.ok) {
        const d = await res.json();
        setSymptoms(d.symptoms ?? "");
        setPrescription(d.prescription ?? "");
        setNotes(d.notes ?? "");
      }
      setLoaded(true);
    });
  }, [appointment.id]);

  const handleSave = async () => {
    if (!symptoms.trim()) {
      setError("Symptoms are required");
      return;
    }
    if (!prescription.trim()) {
      setError("Prescription is required");
      return;
    }
    setError("");
    setSaving(true);
    try {
      const res = await api.saveConsultation({
        appointment_id: appointment.id,
        doctor_id: appointment.doctor_id,
        patient_id: appointment.patient_id,
        clinic_id: appointment.clinic_id,
        symptoms: symptoms.trim(),
        prescription: prescription.trim(),
        notes: notes.trim(),
      });
      if (!res.ok) {
        const d = await res.json();
        setError(d.error ?? "Failed to save");
      } else {
        onClose();
      }
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Visit Pad</h2>
          <button className="modal-close" onClick={onClose}>
            ✕
          </button>
        </div>

        <div className="modal-meta">
          <span>
            <strong>Appointment&nbsp;#</strong>
            {appointment.id}
          </span>
          <span>
            <strong>Doctor&nbsp;</strong>
            {doctorName}
          </span>
          <span>
            <strong>Patient&nbsp;</strong>
            {patientName}
          </span>
          {readOnly && <span className="modal-readonly-tag">Read-only</span>}
        </div>

        {!loaded ? (
          <p className="modal-loading">Loading…</p>
        ) : (
          <>
            <div className="modal-field">
              <label>
                Symptoms <span className="required">*</span>
              </label>
              <input
                value={symptoms}
                onChange={(e) => setSymptoms(e.target.value)}
                placeholder="Enter presenting symptoms…"
                readOnly={readOnly}
              />
            </div>

            <div className="modal-field">
              <label>
                Prescription <span className="required">*</span>
              </label>
              <input
                value={prescription}
                onChange={(e) => setPrescription(e.target.value)}
                placeholder="Enter prescribed medications / treatment…"
                readOnly={readOnly}
              />
            </div>

            <div className="modal-field">
              <label>Notes</label>
              <textarea
                rows={4}
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                placeholder="Additional consultation notes…"
                readOnly={readOnly}
              />
            </div>

            {error && <p className="modal-error">{error}</p>}

            <div className="modal-actions">
              {!readOnly && (
                <button onClick={handleSave} disabled={saving}>
                  {saving ? "Saving…" : "Save"}
                </button>
              )}
              <button onClick={onClose}>Close</button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

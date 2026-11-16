
/* =================================
   MODAL BACKDROP
================================= */

import './EventDetails.css'

.event-details-backdrop {
  position: fixed;
  inset: 0;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px 20px;
  background: rgb(0 0 0 / 45%);
}

/* =================================
   MODAL CONTAINER
================================= */

.event-details {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 760px;
  max-height: calc(100dvh - 64px);
  overflow: hidden;
  border: 1px solid #000;
  border-radius: 8px;
  background: #fff;
  color: #000;
  text-align: left;
}

.event-details::before {
  content: none;
}

/* =================================
   HEADER
================================= */

.event-details-header {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  align-items: center;
  gap: 12px;
  flex: 0 0 auto;
  min-height: 64px;
  padding: 10px 16px;
  border-bottom: 1px solid #e5e4e7;
}

.event-details h2 {
  min-width: 0;
  margin: 0;
  font-size: clamp(22px, 4vw, 32px);
  line-height: 1.15;
  color: #000;
  overflow-wrap: anywhere;
}

/* =================================
   MODAL CONTENT
================================= */

.event-details-content {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  padding: 24px 40px 32px;
  overflow-wrap: anywhere;
}

.event-details-editor-view {
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  padding: 24px 40px 32px;
}

.event-details-editor-view .card {
  padding: 0;
  border: 0;
}

.event-details-club {
  margin-top: 8px;
  font-size: 18px;
  font-weight: 600;
  overflow-wrap: anywhere;
}

.event-details-meta {
  margin-top: 18px;
  font-size: 20px;
  overflow-wrap: anywhere;
}

.event-details-description {
  margin-top: 24px;
  font-size: 18px;
  line-height: 1.6;
  overflow-wrap: anywhere;
}

.event-details-contact {
  margin-top: 20px;
  font-size: 16px;
  overflow-wrap: anywhere;
}

.event-details-contact a {
  overflow-wrap: anywhere;
}

/* =================================
   BUTTON BASE
================================= */

.event-details button {
  font-family: inherit;
  font-weight: 600;
  cursor: pointer;
  transition:
    transform 0.12s ease,
    box-shadow 0.12s ease,
    background 0.12s ease;
}

/* Close button */

.event-details-close {
  padding: 9px 14px;
  border: 1px solid #cbd5e1;
  border-bottom: 3px solid #94a3b8;
  border-radius: 7px;
  background: #f8fafc;
  color: #334155;
  font-size: 14px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.event-details-close:hover {
  background: #e2e8f0;
  transform: translateY(-1px);
  box-shadow: 0 4px 7px rgba(0, 0, 0, 0.12);
}

.event-details-close:active {
  transform: translateY(2px);
  border-bottom-width: 1px;
  box-shadow: none;
}

/* Edit button */

.event-details-edit-toggle {
  padding: 9px 14px;
  border: 1px solid #cbd5e1;
  border-bottom: 3px solid #94a3b8;
  border-radius: 7px;
  background: #f8fafc;
  color: #334155;
  font-size: 14px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.event-details-edit-toggle:hover {
  background: #e2e8f0;
  transform: translateY(-1px);
  box-shadow: 0 4px 7px rgba(0, 0, 0, 0.12);
}

.event-details-edit-toggle:active {
  transform: translateY(2px);
  border-bottom-width: 1px;
  box-shadow: none;
}

/* RSVP button sizing */

.event-details-rsvp {
  min-width: 160px;
  min-height: 44px;
}

/* Cancel registration - destructive action */

.event-details-rsvp-cancel {
  background: #dc2626;
  border-color: #b91c1c;
  border-bottom-color: #991b1b;
}

.event-details-rsvp-cancel:hover:not(:disabled) {
  background: #b91c1c;
  border-color: #991b1b;
  border-bottom-color: #7f1d1d;
}

/* Keyboard focus */

.event-details button:focus-visible {
  outline: 3px solid rgba(37, 99, 235, 0.35);
  outline-offset: 3px;
}

/* =================================
   FOOTER ACTIONS
================================= */

.event-details-actions {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
  gap: 12px;
  padding: 12px 40px 20px;
  border-top: 1px solid #e5e4e7;
  background: #fff;
}

/* =================================
   MOBILE RESPONSIVENESS
================================= */

@media (max-width: 480px) {
  .event-details-backdrop {
    padding: 12px;
  }

  .event-details-content,
  .event-details-actions,
  .event-details-editor-view {
    padding-left: 20px;
    padding-right: 20px;
  }

  .event-details-header {
    gap: 8px;
    padding-left: 12px;
    padding-right: 12px;
  }

  .event-details-actions {
    flex-wrap: wrap;
  }

  .event-details h2 {
    font-size: 22px;
  }
}

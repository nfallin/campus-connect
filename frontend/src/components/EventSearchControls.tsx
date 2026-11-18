import { useEffect, useState } from 'react'
import type { Event } from '../types'
import {
  defaultEventCriteria,
  type EventCriteria,
  type EventSortField,
} from '../lib/eventFilters'
import './EventSearchControls.css'

type EventSearchControlsProps = {
  events: Event[]
  showRsvpFilter?: boolean
  onCriteriaChange: (criteria: EventCriteria) => void
  resultCount: number
}

function uniqueSorted(values: string[]) {
  return [...new Set(values.map((value) => value.trim()).filter(Boolean))]
    .sort((a, b) => a.localeCompare(b, undefined, { sensitivity: 'base', numeric: true }))
}

function EventSearchControls({
  events,
  showRsvpFilter = true,
  onCriteriaChange,
  resultCount,
}: EventSearchControlsProps) {
  const [criteria, setCriteria] = useState<EventCriteria>(defaultEventCriteria)
  const clubs = uniqueSorted(events.map((event) => event.clubName || 'Club not specified'))
  const locations = uniqueSorted(events.map((event) => event.location))
  const activeFilterCount = [
    criteria.search,
    criteria.club !== 'all',
    criteria.location !== 'all',
    criteria.startDate,
    criteria.endDate,
    criteria.approval !== 'all',
    showRsvpFilter && criteria.rsvp !== 'all',
  ].filter(Boolean).length

  useEffect(() => {
    onCriteriaChange(criteria)
  }, [criteria, onCriteriaChange])

  function updateCriteria<Key extends keyof EventCriteria>(
    key: Key,
    value: EventCriteria[Key],
  ) {
    setCriteria((current) => ({ ...current, [key]: value }))
  }

  function clearFilters() {
    setCriteria({
      ...defaultEventCriteria,
      rsvp: showRsvpFilter ? 'all' : 'registered',
    })
  }

  return (
    <section className="event-search" aria-label="Search and filter events">
      <div className="event-search-row">
        <label className="event-search-field event-search-text">
          <span>Search</span>
          <input
            type="search"
            value={criteria.search}
            onChange={(event) => updateCriteria('search', event.target.value)}
            placeholder="Title, club, location, description..."
          />
        </label>
      </div>

      <details className="event-search-advanced">
        <summary>
          Filters &amp; sort
          {activeFilterCount > 0 ? ` (${activeFilterCount})` : ''}
        </summary>

        <div className="event-filter-row">
          <label className="event-search-field">
            <span>Club</span>
            <select
              value={criteria.club}
              onChange={(event) => updateCriteria('club', event.target.value)}
            >
              <option value="all">All clubs</option>
              {clubs.map((club) => <option key={club} value={club}>{club}</option>)}
            </select>
          </label>

          <label className="event-search-field">
            <span>Location</span>
            <select
              value={criteria.location}
              onChange={(event) => updateCriteria('location', event.target.value)}
            >
              <option value="all">All locations</option>
              {locations.map((location) => <option key={location} value={location}>{location}</option>)}
            </select>
          </label>

          <label className="event-search-field">
            <span>Status</span>
            <select
              value={criteria.approval}
              onChange={(event) => updateCriteria(
                'approval',
                event.target.value as EventCriteria['approval'],
              )}
            >
              <option value="all">Any status</option>
              <option value="approved">Approved</option>
              <option value="pending">Pending</option>
            </select>
          </label>

          {showRsvpFilter && (
            <label className="event-search-field">
              <span>RSVP</span>
              <select
                value={criteria.rsvp}
                onChange={(event) => updateCriteria(
                  'rsvp',
                  event.target.value as EventCriteria['rsvp'],
                )}
              >
                <option value="all">All RSVP statuses</option>
                <option value="registered">Registered</option>
                <option value="notRegistered">Not registered</option>
              </select>
            </label>
          )}

          <label className="event-search-field">
            <span>From date</span>
            <input
              type="date"
              value={criteria.startDate}
              onChange={(event) => updateCriteria('startDate', event.target.value)}
            />
          </label>

          <label className="event-search-field">
            <span>Through date</span>
            <input
              type="date"
              value={criteria.endDate}
              min={criteria.startDate || undefined}
              onChange={(event) => updateCriteria('endDate', event.target.value)}
            />
          </label>

          <label className="event-search-field">
            <span>Sort by</span>
            <select
              value={criteria.sortField}
              onChange={(event) => updateCriteria('sortField', event.target.value as EventSortField)}
            >
              <option value="date">Start date</option>
              <option value="endDate">End date</option>
              <option value="title">Title</option>
              <option value="club">Club</option>
              <option value="location">Location</option>
              <option value="approval">Approval status</option>
              {showRsvpFilter && <option value="rsvp">RSVP status</option>}
            </select>
          </label>

          <label className="event-search-field">
            <span>Direction</span>
            <button
              type="button"
              className="button button-secondary"
              onClick={() => updateCriteria(
                'sortDirection',
                criteria.sortDirection === 'asc' ? 'desc' : 'asc',
              )}
              aria-label={`Sort ${criteria.sortDirection === 'asc' ? 'descending' : 'ascending'}`}
            >
              {criteria.sortDirection === 'asc' ? 'Ascending' : 'Descending'}
            </button>
          </label>
        </div>
      </details>

      <div className="event-search-summary">
        <span>{resultCount} {resultCount === 1 ? 'event' : 'events'}</span>
        {activeFilterCount > 0 && (
          <button type="button" className="event-clear-button" onClick={clearFilters}>
            Clear filters
          </button>
        )}
      </div>
    </section>
  )
}

export default EventSearchControls

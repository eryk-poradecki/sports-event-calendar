/**
 * @typedef {Object} EventDetails
 * @property {number} id
 * @property {string} sport_name
 * @property {string} competition_name
 * @property {string} venue_name
 * @property {string} home_team_name
 * @property {string} away_team_name
 * @property {string} start_time
 * @property {"scheduled" | "finished" | "cancelled"} status
 * @property {?number} home_score
 * @property {?number} away_score
 * @property {?string} description
 * @property {boolean} is_neutral_venue
 * @property {string} created_at
 * @property {string} updated_at
 */

/** @returns {Promise<EventDetails[]>} */
async function loadEvents() {
    const response = await fetch("/api/v1/events");

    if (!response.ok) {
        throw new Error("failed to fetch events");
    }
    return await response.json();
}

function groupEventsByDate(events) {
    const grouped = {};

    for (const event of events) {
        const dateKey = new Date(event.start_time).toISOString().split("T")[0];

        if (!grouped[dateKey]) {
            grouped[dateKey] = [];
        }

        grouped[dateKey].push(event);
    }

    return grouped;
}

function formatDate(dateKey) {
    const date = new Date(dateKey);

    return date.toLocaleDateString("en-US", {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric",
    });
}

function formatTime(startTime) {
    const date = new Date(startTime);

    return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    });
}

/** @param {EventDetails} event */
function createEventRow(event) {
    const row = document.createElement("div");
    row.className = "event-row";

    const left = document.createElement("div");
    left.className = "event-meta";

    const title = document.createElement("div");
    title.className = "event-title";
    title.textContent = event.competition_name || event.sport_name;

    left.appendChild(title);

    if (event.competition_name) {
        const subtitle = document.createElement("div");
        subtitle.className = "event-subtitle";
        subtitle.textContent = event.sport_name;
        left.appendChild(subtitle);
    }

    const center = document.createElement("div");
    center.className = "event-match";
    center.textContent = `${event.home_team_name} vs ${event.away_team_name}`;

    const right = document.createElement("div");
    right.className = "event-time";
    right.textContent = formatTime(event.start_time);

    row.appendChild(left);
    row.appendChild(center);
    row.appendChild(right);

    return row;
}

/** @param {EventDetails} events */
function renderEvents(events) {
    const eventsList = document.getElementById("events-list");
    eventsList.innerHTML = "";

    if (!events || events.length === 0) {
        eventsList.textContent = "No events found";
        return;
    }

    const groupedEvents = groupEventsByDate(events);

    for (const dateKey of Object.keys(groupedEvents).sort()) {
        const section = document.createElement("section");
        section.className = "events-date-group";

        const heading = document.createElement("h2");
        heading.className = "events-date-heading";
        heading.textContent = formatDate(dateKey);

        section.appendChild(heading);

        for (const event of groupedEvents[dateKey]) {
            const row = createEventRow(event);
            section.appendChild(row);
        }

        eventsList.appendChild(section);
    }
}

function showError(message) {
    const eventsList = document.getElementById("events-list");
    eventsList.innerHTML = `<p>${message}</p>`;
}

async function init() {
    try {
        const events = await loadEvents();
        renderEvents(events);
    } catch (error) {
        showError(error.message);
    }
}

window.addEventListener('DOMContentLoaded', init);
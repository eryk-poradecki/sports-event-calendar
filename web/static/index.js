/**
 * @typedef {Object} EventListItem
 * @property {number} id
 * @property {string} sport_name
 * @property {string} competition_name
 * @property {string} venue_name
 * @property {string} home_team_name
 * @property {string} away_team_name
 * @property {string} start_time
 * @property {"scheduled" | "finished" | "cancelled"} status
 */

let currentPage = 1;
let currentFilters = {
  sport: "",
  dateFrom: "",
  dateTo: "",
};

const pageSize = 10;

/** @returns {Promise<EventListItem[]>} */
async function loadEvents(page = 1) {
  const params = new URLSearchParams({
    page: String(page),
    page_size: String(pageSize),
  });

  if (currentFilters.sport) {
    params.set("sport", currentFilters.sport);
  }

  if (currentFilters.dateFrom) {
    params.set("date_from", currentFilters.dateFrom);
  }

  if (currentFilters.dateTo) {
    params.set("date_to", currentFilters.dateTo);
  }

  console.log(params);
  const response = await fetch(`/api/v1/events?${params.toString()}`);

  if (!response.ok) {
    throw new Error("failed to fetch events");
  }
  return await response.json();
}

async function loadSports() {
  const response = await fetch(`/api/v1/sports`);

  if (!response.ok) {
    throw new Error("failed to fetch sports");
  }

  return await response.json();
}

async function loadTeams() {
  const response = await fetch(`/api/v1/teams`);
  if (!response.ok) {
    throw new Error("failed to fetch teams");
  }
  return await response.json();
}

async function loadCompetitions() {
  const response = await fetch(`/api/v1/competitions`);
  if (!response.ok) {
    throw new Error("failed to fetch competitions");
  }
  return await response.json();
}

async function loadVenues() {
  const response = await fetch(`/api/v1/venues`);
  if (!response.ok) {
    throw new Error("failed to fetch venues");
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

/** @param {EventListItem} eventListItem */
function createEventRow(eventListItem) {
  const row = document.createElement("div");
  row.className = "event-row";
  row.addEventListener("click", async (event) => {
    window.location.href = `/events/${eventListItem.id}`;
  });

  const left = document.createElement("div");
  left.className = "event-meta";

  const title = document.createElement("div");
  title.className = "event-title";
  title.textContent =
    eventListItem.competition_name || eventListItem.sport_name;

  left.appendChild(title);

  if (eventListItem.competition_name) {
    const subtitle = document.createElement("div");
    subtitle.className = "event-subtitle";
    subtitle.textContent = eventListItem.sport_name;
    left.appendChild(subtitle);
  }

  const center = document.createElement("div");
  center.className = "event-match";
  if (eventListItem.venue_name !== null) {
    center.textContent = `${eventListItem.home_team_name} vs ${eventListItem.away_team_name} at ${eventListItem.venue_name}`;
  } else {
    center.textContent = `${eventListItem.home_team_name} vs ${eventListItem.away_team_name}`;
  }

  const right = document.createElement("div");
  right.className = "event-time";
  right.textContent = formatTime(eventListItem.start_time);

  row.appendChild(left);
  row.appendChild(center);
  row.appendChild(right);

  return row;
}

/** @param {EventListItem} events */
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

function renderPagination(page, totalPages) {
  const pagination = document.getElementById("pagination");
  pagination.innerHTML = "";

  if (totalPages <= 1) {
    return;
  }

  const prevButton = document.createElement("button");
  prevButton.textContent = "Previous";
  prevButton.disabled = page <= 1;
  prevButton.addEventListener("click", (event) => {
    refreshEvents(page - 1);
  });

  const info = document.createElement("span");
  info.textContent = `Page ${page} of ${totalPages}`;

  const nextButton = document.createElement("button");
  nextButton.textContent = "Next";
  nextButton.disabled = page >= totalPages;
  nextButton.addEventListener("click", (event) => {
    refreshEvents(page + 1);
  });

  pagination.appendChild(prevButton);
  pagination.appendChild(info);
  pagination.appendChild(nextButton);
}

function showError(message) {
  const eventsList = document.getElementById("events-list");
  eventsList.innerHTML = `<p>${message}</p>`;
}

async function refreshEvents(page = 1) {
  try {
    const data = await loadEvents(page);
    currentPage = data.page;

    renderEvents(data.items);
    renderPagination(data.page, data.total_pages);
  } catch (error) {
    showError(error.message);
  }
}

function renderFilters(sports) {
  const filtersSection = document.getElementById("sport-filter-section");
  filtersSection.innerHTML = "";

  const form = document.createElement("form");
  form.className = "filters-form";

  const sportLabel = document.createElement("label");
  sportLabel.setAttribute("for", "sport-filter");
  sportLabel.textContent = "Sport";

  const sportSelect = document.createElement("select");
  sportSelect.id = "sport-filter";
  sportSelect.name = "sport";

  const defaultOption = document.createElement("option");
  defaultOption.value = "";
  defaultOption.textContent = "All sports";
  sportSelect.appendChild(defaultOption);

  for (const sport of sports) {
    const option = document.createElement("option");
    option.value = sport.slug;
    option.textContent = sport.name;
    sportSelect.appendChild(option);
  }

  sportSelect.value = currentFilters.sport;

  sportSelect.addEventListener("change", async (event) => {
    currentFilters.sport = event.target.value;
    currentPage = 1;
    await refreshEvents(currentPage);
  });

  form.appendChild(sportLabel);
  form.appendChild(sportSelect);

  filtersSection.appendChild(form);

  const dateFrom = document.getElementById("date-from");
  const dateTo = document.getElementById("date-to");

  dateFrom.addEventListener("change", async (event) => {
    currentFilters.dateFrom = event.target.value;

    if (event.target.value) {
      dateTo.min = event.target.value;

      if (dateTo.value && dateTo.value < event.target.value) {
        dateTo.value = "";
        currentFilters.dateTo = "";
      }
    } else {
      dateTo.removeAttribute("min");
    }

    currentPage = 1;
    await refreshEvents(currentPage);
  });

  dateTo.addEventListener("change", async (event) => {
    currentFilters.dateTo = event.target.value;
    currentPage = 1;
    await refreshEvents(currentPage);
  });
}

async function clearFilters() {
  const dateFrom = document.getElementById("date-from");
  const dateTo = document.getElementById("date-to");
  const sportSelect = document.getElementById("sport-filter");

  dateFrom.value = "";
  dateTo.value = "";
  sportSelect.value = "";

  currentFilters.sport = "";
  currentFilters.dateFrom = "";
  currentFilters.dateTo = "";

  currentPage = 1;
  await refreshEvents(currentPage);
}

function populateSelect(
  selectElement,
  items,
  valueKey,
  labelKey,
  placeholder = null,
) {
  selectElement.innerHTML = "";

  if (placeholder !== null) {
    const placeholderOption = document.createElement("option");
    placeholderOption.value = "";
    placeholderOption.textContent = placeholder;
    selectElement.appendChild(placeholderOption);
  }

  for (const item of items) {
    const option = document.createElement("option");
    option.value = item[valueKey];
    option.textContent = item[labelKey];
    selectElement.appendChild(option);
  }
}

function toggleScoreFields() {
  const statusSelect = document.getElementById("status");
  const scoreFields = document.getElementById("score-fields");

  if (statusSelect.value === "finished") {
    scoreFields.hidden = false;
  } else {
    scoreFields.hidden = true;
    document.getElementById("home-score").value = "";
    document.getElementById("away-score").value = "";
  }
}

function formatDateTimeLocalToISOString(value) {
  return new Date(value).toISOString();
}

async function submitCreateEventForm(event) {
  event.preventDefault();

  const message = document.getElementById("create-event-message");
  message.textContent = "";

  const sportId = document.getElementById("sport-id").value;
  const competitionId = document.getElementById("competition-id").value;
  const venueId = document.getElementById("venue-id").value;
  const homeTeamId = document.getElementById("home-team-id").value;
  const awayTeamId = document.getElementById("away-team-id").value;
  const startTime = document.getElementById("start-time").value;
  const status = document.getElementById("status").value;
  const isNeuralVenue = document.getElementById("is-neutral-venue").checked;
  const homeScore = document.getElementById("home-score").value;
  const awayScore = document.getElementById("away-score").value;
  const description = document.getElementById("description").value;

  const payload = {
    _sport_id: Number(sportId),
    _home_team_id: Number(homeTeamId),
    _away_team_id: Number(awayTeamId),
    start_time: formatDateTimeLocalToISOString(startTime),
    status: status,
    is_neutral_venue: isNeuralVenue,
  };

  if (competitionId) {
    payload._competition_id = Number(competitionId);
  }

  if (venueId) {
    payload._venue_id = Number(venueId);
  }

  if (status === "finished") {
    if (homeScore !== "") {
      payload.home_score = Number(homeScore);
    }
    if (awayScore !== "") {
      payload.away_score = Number(awayScore);
    }
    }
  if (description) {
      payload.description = description;
  }

  try {
      const response = await fetch("/api/v1/events", {
          method: "POST",
          headers: {
              "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
      });

      if (!response.ok) {
          const errorText = await response.text();
          throw new Error(errorText || "failed to create event");
      }

      message.textContent = "Event created successfully.";
      document.getElementById("create-event-form").reset();
      toggleScoreFields();
      const form = document.getElementById("create-event-form");
      setTimeout(() => {form.hidden = true;}, 3000);
      currentPage = 1;
      await refreshEvents(currentPage);
  } catch (error) {
      message.textContent = error;
  }
}

async function initCreateEventForm() {
    const toggleButton = document.getElementById("toggle-event-form-btn");
    const form = document.getElementById("create-event-form");
    const statusSelect = document.getElementById("status");

    toggleButton.addEventListener("click", () => {
        form.hidden = !form.hidden;
    });

    statusSelect.addEventListener("change", toggleScoreFields);

    const [sports, teams, competitions, venues] = await Promise.all([
        loadSports(),
        loadTeams(),
        loadCompetitions(),
        loadVenues(),
    ]);

    populateSelect(document.getElementById("sport-id"), sports, "id", "name");
    populateSelect(document.getElementById("home-team-id"), teams, "id", "name");
    populateSelect(document.getElementById("away-team-id"), teams, "id", "name");
    populateSelect(document.getElementById("competition-id"), competitions, "id", "name", "No competition");
    populateSelect(document.getElementById("venue-id"), venues, "id", "name", "No venue");

    form.addEventListener("submit", submitCreateEventForm);
    toggleScoreFields();
}

async function initFilters() {
  try {
    const sports = await loadSports();
    renderFilters(sports);
  } catch (error) {
    console.error(error);
  }
}

async function init() {
  await initFilters();
  await initCreateEventForm();
  await refreshEvents(1);
}

window.addEventListener("DOMContentLoaded", init);
window.clearFilters = clearFilters;


/**
 * go_vibe - Weather Forecast & Analytics Dashboard
 * Frontend & Backend Integration Script
 */

(function () {
  'use strict';

  // 1. API Configuration & Environment Detection
  const API_BASE = (window.location.protocol === 'file:' || !window.location.port)
    ? 'http://localhost:8080'
    : '';

  // 2. Global Application State
  const state = {
    currentLocationId: null,
    currentLocation: {
      id: '00000000-0000-0000-0000-000000000001',
      city: 'Los Angeles',
      state: 'CA',
      country: 'USA',
      image_url: 'https://images.unsplash.com/photo-1580655653885-65763b2597d0?auto=format&fit=crop&w=800&q=80'
    },
    unit: localStorage.getItem('temperature_unit') || 'Celsius', // 'Celsius' | 'Fahrenheit'
    token: localStorage.getItem('jwt_token') || null,
    user: null,
    savedLocations: [],
    chartPeriod: 'today', // 'today' | 'week'
    dashboardData: null,
    compareCity1: 'Los Angeles',
    compareCity2: 'Jakarta',
    isCompareMode: false,
    comparisonData: null,
    calendarCity: 'Los Angeles',
    calendarYear: 2026,
    calendarMonth: 10,
    calendarData: null,
    selectedDayIndex: 0,
    isCalendarMode: false
  };

  // 3. API Request Helper
  async function apiRequest(endpoint, options = {}) {
    const headers = {
      'Accept': 'application/json',
      ...(options.headers || {})
    };

    if (state.token) {
      headers['Authorization'] = `Bearer ${state.token}`;
    }

    if (options.body && typeof options.body === 'object' && !(options.body instanceof FormData)) {
      headers['Content-Type'] = 'application/json';
      options.body = JSON.stringify(options.body);
    }

    try {
      const url = `${API_BASE}${endpoint}`;
      const response = await fetch(url, { ...options, headers });
      const data = await response.json();

      if (!response.ok || !data.success) {
        throw new Error(data.message || `HTTP Error ${response.status}`);
      }

      return data;
    } catch (err) {
      console.warn(`[API] Error on ${endpoint}:`, err.message);
      throw err;
    }
  }

  // 4. Toast Notification System
  function showToast(message, type = 'info') {
    const container = document.getElementById('toastContainer');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;
    toast.innerHTML = `<span>${message}</span>`;
    container.appendChild(toast);

    setTimeout(() => {
      toast.style.opacity = '0';
      toast.style.transform = 'translateY(10px)';
      toast.style.transition = 'all 0.25s ease';
      setTimeout(() => toast.remove(), 250);
    }, 3200);
  }

  // 5. SVG Weather Icons Library
  function getWeatherIconSvg(condition = '', isHero = false) {
    const cond = condition.toLowerCase();

    if (cond.includes('sun') || cond.includes('clear')) {
      return `
        <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
          <circle cx="32" cy="32" r="14" fill="#f8b42c"></circle>
          <g stroke="#f8b42c" stroke-width="3" stroke-linecap="round">
            <line x1="32" y1="8" x2="32" y2="13"></line>
            <line x1="32" y1="51" x2="32" y2="56"></line>
            <line x1="8" y1="32" x2="13" y2="32"></line>
            <line x1="51" y1="32" x2="56" y2="32"></line>
            <line x1="15" y1="15" x2="18.5" y2="18.5"></line>
            <line x1="45.5" y1="45.5" x2="49" y2="49"></line>
            <line x1="15" y1="49" x2="18.5" y2="45.5"></line>
            <line x1="45.5" y1="18.5" x2="49" y2="15"></line>
          </g>
        </svg>
      `;
    }

    if (cond.includes('snow') || cond.includes('ice')) {
      return `
        <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
          <path d="M18 38h28a10 10 0 0 0 1.2-19.9 14 14 0 0 0-27.4 3.7A8 8 0 0 0 18 38z" fill="#ffffff" fill-opacity="0.95"></path>
          <circle cx="24" cy="46" r="2" fill="#8ca5c2"></circle>
          <circle cx="32" cy="46" r="2" fill="#8ca5c2"></circle>
          <circle cx="40" cy="46" r="2" fill="#8ca5c2"></circle>
          <circle cx="28" cy="52" r="1.5" fill="#8ca5c2"></circle>
          <circle cx="36" cy="52" r="1.5" fill="#8ca5c2"></circle>
        </svg>
      `;
    }

    if (cond.includes('thunder') || cond.includes('storm')) {
      return `
        <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
          <path d="M18 36h28a10 10 0 0 0 1.2-19.9 14 14 0 0 0-27.4 3.7A8 8 0 0 0 18 36z" fill="#4d5382"></path>
          <polygon points="30,38 24,48 31,48 27,58 38,46 31,46" fill="#f8b42c"></polygon>
        </svg>
      `;
    }

    if (cond.includes('rain') || cond.includes('drizzle') || cond.includes('shower')) {
      return `
        <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
          <path d="M18 38h28a10 10 0 0 0 1.2-19.9 14 14 0 0 0-27.4 3.7A8 8 0 0 0 18 38z" fill="#ffffff" fill-opacity="0.95"></path>
          <g stroke="#7e9ebf" stroke-width="2.5" stroke-linecap="round">
            <line x1="22" y1="44" x2="19" y2="52"></line>
            <line x1="32" y1="44" x2="29" y2="52"></line>
            <line x1="42" y1="44" x2="39" y2="52"></line>
          </g>
        </svg>
      `;
    }

    if (cond.includes('cloud') || cond.includes('overcast')) {
      return `
        <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
          <circle cx="26" cy="24" r="10" fill="#f8b42c"></circle>
          <path d="M18 42h28a10 10 0 0 0 1.2-19.9 14 14 0 0 0-27.4 3.7A8 8 0 0 0 18 42z" fill="#ffffff" fill-opacity="0.95"></path>
        </svg>
      `;
    }

    // Default icon
    return `
      <svg class="${isHero ? 'weather-illustration' : ''}" viewBox="0 0 64 64" fill="none">
        <circle cx="26" cy="24" r="10" fill="#f8b42c"></circle>
        <path d="M18 42h28a10 10 0 0 0 1.2-19.9 14 14 0 0 0-27.4 3.7A8 8 0 0 0 18 42z" fill="#ffffff" fill-opacity="0.95"></path>
      </svg>
    `;
  }

  // 6. SVG Smooth Bézier Curve Generator
  function renderSvgChart(points, labels, unitSymbol) {
    const svg = document.getElementById('hourlyChartSvg');
    if (!svg || !points || points.length === 0) return;

    const width = 540;
    const height = 110;
    const paddingX = 35;
    const paddingTop = 26;
    const paddingBottom = 22;

    const min = Math.min(...points);
    const max = Math.max(...points);
    const range = (max - min) === 0 ? 1 : (max - min);

    const stepX = (width - paddingX * 2) / (points.length - 1);
    const coords = points.map((val, idx) => {
      const x = paddingX + idx * stepX;
      const normalized = (val - min) / range;
      const y = (height - paddingBottom) - (normalized * (height - paddingTop - paddingBottom));
      return { x, y, val };
    });

    // Smooth Bézier curve generator
    function buildSmoothPath(pts) {
      if (pts.length < 2) return '';
      let d = `M ${pts[0].x} ${pts[0].y}`;
      for (let i = 0; i < pts.length - 1; i++) {
        const p0 = pts[i === 0 ? i : i - 1];
        const p1 = pts[i];
        const p2 = pts[i + 1];
        const p3 = pts[i + 2] || p2;

        const cp1x = p1.x + (p2.x - p0.x) / 6;
        const cp1y = p1.y + (p2.y - p0.y) / 6;
        const cp2x = p2.x - (p3.x - p1.x) / 6;
        const cp2y = p2.y - (p3.y - p1.y) / 6;

        d += ` C ${cp1x.toFixed(1)} ${cp1y.toFixed(1)}, ${cp2x.toFixed(1)} ${cp2y.toFixed(1)}, ${p2.x.toFixed(1)} ${p2.y.toFixed(1)}`;
      }
      return d;
    }

    const curveD = buildSmoothPath(coords);
    const areaD = `${curveD} L ${coords[coords.length - 1].x} ${height} L ${coords[0].x} ${height} Z`;

    let markersHtml = '';
    coords.forEach((pt) => {
      markersHtml += `
        <g class="chart-point-group">
          <circle cx="${pt.x.toFixed(1)}" cy="${pt.y.toFixed(1)}" r="4" fill="#ffffff" stroke="#8c7cf0" stroke-width="2.5"></circle>
          <text x="${pt.x.toFixed(1)}" y="${(pt.y - 10).toFixed(1)}" text-anchor="middle" font-size="11" font-weight="700" fill="#2c2d4a" font-family="inherit">${Math.round(pt.val)}°</text>
        </g>
      `;
    });

    svg.innerHTML = `
      <defs>
        <linearGradient id="chartGradient" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#8c7cf0" stop-opacity="0.25"/>
          <stop offset="100%" stop-color="#8c7cf0" stop-opacity="0.0"/>
        </linearGradient>
      </defs>
      <path d="${areaD}" fill="url(#chartGradient)"></path>
      <path d="${curveD}" fill="none" stroke="#8c7cf0" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"></path>
      ${markersHtml}
    `;

    // Render labels
    const labelsEl = document.getElementById('chartTimeLabels');
    if (labelsEl && labels && labels.length > 0) {
      labelsEl.innerHTML = labels.map(l => `<span>${l}</span>`).join('');
    }
  }

  // 7. Render Chance of Rain Section
  function renderChanceOfRain(rainList) {
    const listEl = document.getElementById('rainBarsList');
    if (!listEl || !rainList || rainList.length === 0) return;

    let maxChance = -1;
    rainList.forEach(item => {
      if (item.percentage > maxChance) maxChance = item.percentage;
    });

    listEl.innerHTML = rainList.map(item => {
      const isPeak = item.percentage === maxChance && maxChance > 0;
      return `
        <div class="rain-bar-row ${isPeak ? 'active-highlight-row' : ''}">
          <span class="rain-time">${item.time}</span>
          <div class="bar-track" title="${item.percentage}% chance (${item.intensity})">
            <div class="bar-fill ${isPeak ? 'dark-navy-bar' : 'purple-bar'}" style="width: ${Math.min(100, item.percentage)}%;"></div>
          </div>
        </div>
      `;
    }).join('');
  }

  // 8. Render 3 Days Forecast Section
  function renderForecast(forecastList) {
    const listEl = document.getElementById('forecastCardsList');
    if (!listEl || !forecastList || forecastList.length === 0) return;

    listEl.innerHTML = forecastList.map((item, idx) => {
      const isLeft = idx % 2 === 0;
      const highTemp = Math.round(item.temp_high);
      const lowTemp = Math.round(item.temp_low);
      const iconSvg = getWeatherIconSvg(item.condition);

      const pillHtml = `
        <div class="temp-pill">
          <span class="arrow-up">↑</span> ${highTemp}°
          <span class="arrow-down">↓</span> ${lowTemp}°
        </div>
      `;

      const infoHtml = `
        <div class="forecast-info-group">
          <div class="forecast-icon">
            ${iconSvg}
          </div>
          <div class="forecast-text">
            <div class="forecast-day">${item.day}</div>
            <div class="forecast-cond">${item.condition}</div>
          </div>
        </div>
      `;

      return `
        <div class="forecast-card-row ${isLeft ? 'pill-left' : 'pill-right'}">
          ${isLeft ? pillHtml + infoHtml : infoHtml + pillHtml}
        </div>
      `;
    }).join('');
  }

  // 9. Render Main Dashboard Header & Highlights
  function renderDashboard(data) {
    state.dashboardData = data;
    const loc = data.location;
    const cur = data.current_weather;
    const high = data.highlights;

    // Location
    const locNameEl = document.getElementById('locationTitle');
    if (locNameEl && loc) {
      const stateStr = loc.state ? `, ${loc.state}` : '';
      locNameEl.textContent = `${loc.city}${stateStr}, ${loc.country}`;
      state.calendarCity = loc.city;
    }

    // Hero background photo
    const heroCard = document.getElementById('weatherHeroCard');
    if (heroCard && loc && loc.image_url) {
      heroCard.style.backgroundImage = `url('${loc.image_url}')`;
    }

    // Current temperature & condition
    const tempEl = document.getElementById('currentTemp');
    if (tempEl && cur) {
      tempEl.textContent = Math.round(cur.temperature);
    }

    const unitLabel = document.getElementById('tempUnitLabel');
    if (unitLabel) {
      unitLabel.textContent = state.unit === 'Fahrenheit' ? '°F' : '°C';
    }

    const dtEl = document.getElementById('currentDateTime');
    if (dtEl && cur) {
      dtEl.textContent = cur.date_time || 'Today';
    }

    const condEl = document.getElementById('currentCondition');
    if (condEl && cur) {
      condEl.textContent = cur.condition;
    }

    const heroIconEl = document.getElementById('heroWeatherIcon');
    if (heroIconEl && cur) {
      heroIconEl.innerHTML = getWeatherIconSvg(cur.condition, true);
    }

    // Highlights
    if (high) {
      const precipEl = document.getElementById('highlightPrecip');
      if (precipEl) precipEl.textContent = high.precipitation || '0%';

      const humidEl = document.getElementById('highlightHumidity');
      if (humidEl) humidEl.textContent = high.humidity || '0%';

      const windEl = document.getElementById('highlightWind');
      if (windEl) windEl.textContent = high.wind_speed || '0 km/h';

      const sunriseEl = document.getElementById('highlightSunrise');
      if (sunriseEl) sunriseEl.textContent = high.sunrise || '6:00 am';

      const sunsetEl = document.getElementById('highlightSunset');
      if (sunsetEl) sunsetEl.textContent = high.sunset || '6:00 pm';
    }

    // Hourly Chart
    if (data.hourly_temperature && data.hourly_temperature.length > 0) {
      const points = data.hourly_temperature.map(h => h.temp);
      const labels = data.hourly_temperature.map(h => h.time);
      renderSvgChart(points, labels, state.unit === 'Fahrenheit' ? '°F' : '°C');
    }

    // Chance of Rain
    if (data.chance_of_rain) {
      renderChanceOfRain(data.chance_of_rain);
    }

    // 3 Days Forecast
    if (data.three_days_forecast) {
      renderForecast(data.three_days_forecast);
    }
  }

  // 10. Core Data Loading Function
  async function loadDashboard(locationId = null) {
    try {
      const unitParam = state.unit.toLowerCase();
      let query = `?unit=${unitParam}`;
      if (locationId) {
        query += `&location_id=${locationId}`;
        state.currentLocationId = locationId;
      }

      const res = await apiRequest(`/api/v1/weather/dashboard${query}`);
      if (res && res.data) {
        renderDashboard(res.data);
      }
    } catch (err) {
      showToast('Gagal memuat data cuaca dari backend: ' + err.message, 'error');
    }
  }

  // 11. Search Locations Logic
  async function performSearch(query, dropdownEl, onSelect) {
    if (!query || query.trim().length === 0) {
      dropdownEl.classList.remove('active');
      dropdownEl.innerHTML = '';
      return;
    }

    try {
      dropdownEl.innerHTML = '<div class="dropdown-item">Mencari kota...</div>';
      dropdownEl.classList.add('active');

      const res = await apiRequest(`/api/v1/locations/search?q=${encodeURIComponent(query.trim())}`);
      const locations = res.data || [];

      if (locations.length === 0) {
        dropdownEl.innerHTML = '<div class="dropdown-item" style="color:#8c92a4;">Tidak ditemukan kota yang cocok</div>';
        return;
      }

      dropdownEl.innerHTML = locations.map(loc => {
        const stateStr = loc.state ? `, ${loc.state}` : '';
        return `
          <div class="dropdown-item" data-id="${loc.id}">
            <div>
              <strong>${loc.city}</strong>
              <span style="font-size:11px; color:#8c92a4; margin-left:4px;">${stateStr}, ${loc.country}</span>
            </div>
            <span style="font-size:11px; color:#8c7cf0; font-weight:600;">Lihat Cuaca &rarr;</span>
          </div>
        `;
      }).join('');

      dropdownEl.querySelectorAll('.dropdown-item[data-id]').forEach(item => {
        item.addEventListener('click', () => {
          const locId = item.getAttribute('data-id');
          const found = locations.find(l => l.id === locId);
          dropdownEl.classList.remove('active');
          if (onSelect) onSelect(found);
        });
      });
    } catch (err) {
      dropdownEl.innerHTML = `<div class="dropdown-item" style="color:#e53e3e;">Pencarian gagal: ${err.message}</div>`;
    }
  }

  // 12. Setup Top Search Bar
  function setupSearchBar() {
    const input = document.getElementById('searchInput');
    const submitBtn = document.getElementById('searchSubmitBtn');
    const dropdown = document.getElementById('searchDropdown');
    if (!input || !dropdown) return;

    let debounceTimer = null;

    input.addEventListener('input', () => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        performSearch(input.value, dropdown, (selectedLoc) => {
          if (selectedLoc) {
            input.value = `${selectedLoc.city}, ${selectedLoc.country}`;
            loadDashboard(selectedLoc.id);
            showToast(`Menampilkan cuaca untuk ${selectedLoc.city}`, 'success');
          }
        });
      }, 250);
    });

    submitBtn?.addEventListener('click', () => {
      if (input.value.trim()) {
        performSearch(input.value, dropdown, (selectedLoc) => {
          if (selectedLoc) {
            input.value = `${selectedLoc.city}, ${selectedLoc.country}`;
            loadDashboard(selectedLoc.id);
          }
        });
      }
    });

    document.addEventListener('click', (e) => {
      if (!input.contains(e.target) && !dropdown.contains(e.target)) {
        dropdown.classList.remove('active');
      }
    });
  }

  // 13. Setup Add City Modal
  function setupCityModal() {
    const modal = document.getElementById('cityModal');
    const addBtn = document.getElementById('addCityBtn');
    const closeBtn = document.getElementById('modalCloseBtn');
    const searchInput = document.getElementById('modalSearchInput');
    const searchBtn = document.getElementById('modalSearchBtn');
    const resultsList = document.getElementById('modalResultsList');

    if (!modal || !addBtn) return;

    const openModal = () => {
      modal.classList.add('active');
      if (searchInput) {
        searchInput.value = '';
        searchInput.focus();
      }
      if (resultsList) resultsList.innerHTML = '';
    };

    const closeModal = () => modal.classList.remove('active');

    addBtn.addEventListener('click', openModal);
    closeBtn?.addEventListener('click', closeModal);
    modal.addEventListener('click', (e) => {
      if (e.target === modal) closeModal();
    });

    const doSearch = () => {
      const q = searchInput?.value.trim();
      if (!q) return;

      resultsList.innerHTML = '<div style="padding:10px; color:#8c92a4; text-align:center;">Mencari kota...</div>';
      apiRequest(`/api/v1/locations/search?q=${encodeURIComponent(q)}`)
        .then(res => {
          const locs = res.data || [];
          if (locs.length === 0) {
            resultsList.innerHTML = '<div style="padding:10px; color:#8c92a4; text-align:center;">Tidak ada kota ditemukan</div>';
            return;
          }

          resultsList.innerHTML = locs.map(loc => `
            <div class="saved-loc-item" data-id="${loc.id}">
              <div>
                <div class="saved-loc-title">${loc.city}</div>
                <div class="saved-loc-subtitle">${loc.state ? loc.state + ', ' : ''}${loc.country}</div>
              </div>
              <div style="display:flex; gap:6px;">
                <button class="btn-secondary select-city-btn" data-id="${loc.id}" style="padding:6px 12px; font-size:11.5px;">Pilih</button>
                ${state.token ? `<button class="btn-secondary save-city-btn" data-id="${loc.id}" style="padding:6px 10px; font-size:11.5px;">⭐ Simpan</button>` : ''}
              </div>
            </div>
          `).join('');

          resultsList.querySelectorAll('.select-city-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
              e.stopPropagation();
              const id = btn.getAttribute('data-id');
              const found = locs.find(l => l.id === id);
              closeModal();
              loadDashboard(id);
              if (found) showToast(`Menampilkan cuaca ${found.city}`, 'success');
            });
          });

          resultsList.querySelectorAll('.save-city-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
              e.stopPropagation();
              const id = btn.getAttribute('data-id');
              try {
                await apiRequest('/api/v1/locations/saved', {
                  method: 'POST',
                  body: { location_id: id }
                });
                showToast('Lokasi berhasil disimpan ke daftar favorit!', 'success');
                btn.textContent = ' Tersimpan';
                btn.disabled = true;
              } catch (err) {
                showToast(err.message, 'error');
              }
            });
          });
        })
        .catch(err => {
          resultsList.innerHTML = `<div style="padding:10px; color:#e53e3e;">${err.message}</div>`;
        });
    };

    searchBtn?.addEventListener('click', doSearch);
    searchInput?.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') doSearch();
    });
  }

  // 14. Setup Preferences & Settings Modal
  function setupSettingsModal() {
    const modal = document.getElementById('settingsModal');
    const navBtn = document.getElementById('settingsNavBtn');
    const closeBtn = document.getElementById('settingsCloseBtn');
    const unitBtns = document.querySelectorAll('.unit-btn');

    if (!modal) return;

    const openModal = () => {
      modal.classList.add('active');
      unitBtns.forEach(b => {
        const u = b.getAttribute('data-unit');
        b.classList.toggle('active', u.toLowerCase() === state.unit.toLowerCase());
      });
    };

    const closeModal = () => modal.classList.remove('active');

    navBtn?.addEventListener('click', (e) => {
      e.preventDefault();
      openModal();
    });

    closeBtn?.addEventListener('click', closeModal);
    modal.addEventListener('click', (e) => {
      if (e.target === modal) closeModal();
    });

    unitBtns.forEach(btn => {
      btn.addEventListener('click', async () => {
        const newUnit = btn.getAttribute('data-unit');
        state.unit = newUnit;
        localStorage.setItem('temperature_unit', newUnit);

        unitBtns.forEach(b => b.classList.toggle('active', b === btn));

        // If user is logged in, sync with backend
        if (state.token) {
          try {
            await apiRequest('/api/v1/users/settings', {
              method: 'PUT',
              body: { temperature_unit: newUnit }
            });
          } catch (e) {
            console.warn('Sync settings warning:', e.message);
          }
        }

        showToast(`Satuan suhu diubah ke ${newUnit}`, 'success');
        if (state.isCompareMode) {
          loadComparison(state.compareCity1, state.compareCity2);
        } else if (state.isCalendarMode) {
          loadCalendar();
        } else {
          loadDashboard(state.currentLocationId);
        }
      });
    });
  }

  // 15. Setup User Auth & Profile Modal
  function setupAuthModal() {
    const modal = document.getElementById('authModal');
    const profileBtn = document.getElementById('profileBtn');
    const avatarImg = document.getElementById('avatarImg');
    const closeBtn = document.getElementById('authCloseBtn');
    const modalBody = document.getElementById('authModalBody');
    const modalTitle = document.getElementById('authModalTitle');

    if (!modal || !modalBody) return;

    const renderAuthContent = () => {
      if (state.token && state.user) {
        // Logged-in profile view
        if (modalTitle) modalTitle.textContent = 'User Profile & Preferences';
        modalBody.innerHTML = `
          <div class="profile-card">
            <div class="profile-info">
              <img src="https://images.unsplash.com/photo-1534528741775-53994a69daeb?auto=format&fit=crop&w=150&q=80" alt="Avatar" class="profile-info-avatar">
              <div>
                <div class="profile-name">${state.user.name || 'Pengguna'}</div>
                <div class="profile-email">${state.user.email}</div>
                <div style="font-size:11.5px; color:#8c7cf0; margin-top:2px; font-weight:600;">Unit Suhu: ${state.user.temperature_unit || state.unit}</div>
              </div>
            </div>

            <div class="profile-actions">
              <button class="btn-secondary" id="viewSavedLocationsBtn">Kota Tersimpan</button>
              <button class="btn-danger" id="logoutBtn">Logout</button>
            </div>
          </div>
        `;

        document.getElementById('logoutBtn')?.addEventListener('click', () => {
          localStorage.removeItem('jwt_token');
          state.token = null;
          state.user = null;
          showToast('Anda telah keluar dari akun.', 'info');
          closeModal();
          loadDashboard(state.currentLocationId);
        });

        document.getElementById('viewSavedLocationsBtn')?.addEventListener('click', () => {
          closeModal();
          openSavedLocationsModal();
        });
      } else {
        // Guest / Login & Register Form
        if (modalTitle) modalTitle.textContent = 'Login or Register';
        modalBody.innerHTML = `
          <div class="auth-tabs">
            <button class="auth-tab-btn active" id="tabLoginBtn">Masuk</button>
            <button class="auth-tab-btn" id="tabRegisterBtn">Daftar Akun</button>
          </div>

          <form id="authForm" class="auth-form">
            <div class="form-field" id="nameField" style="display:none;">
              <label>Nama Lengkap</label>
              <input type="text" id="authNameInput" placeholder="Masukkan nama Anda">
            </div>

            <div class="form-field">
              <label>Email</label>
              <input type="email" id="authEmailInput" placeholder="email@example.com" required>
            </div>

            <div class="form-field">
              <label>Password</label>
              <input type="password" id="authPasswordInput" placeholder="Minimal 6 karakter" required>
            </div>

            <button type="submit" class="auth-submit-btn" id="authSubmitBtn">Masuk ke Akun</button>
          </form>
        `;

        let isRegister = false;
        const tabLogin = document.getElementById('tabLoginBtn');
        const tabRegister = document.getElementById('tabRegisterBtn');
        const nameField = document.getElementById('nameField');
        const submitBtn = document.getElementById('authSubmitBtn');
        const form = document.getElementById('authForm');

        tabLogin?.addEventListener('click', () => {
          isRegister = false;
          tabLogin.classList.add('active');
          tabRegister.classList.remove('active');
          nameField.style.display = 'none';
          submitBtn.textContent = 'Masuk ke Akun';
        });

        tabRegister?.addEventListener('click', () => {
          isRegister = true;
          tabRegister.classList.add('active');
          tabLogin.classList.remove('active');
          nameField.style.display = 'flex';
          submitBtn.textContent = 'Daftar Akun Baru';
        });

        form?.addEventListener('submit', async (e) => {
          e.preventDefault();
          const email = document.getElementById('authEmailInput')?.value.trim();
          const password = document.getElementById('authPasswordInput')?.value.trim();
          const name = document.getElementById('authNameInput')?.value.trim();

          if (!email || !password) {
            showToast('Email dan password wajib diisi', 'error');
            return;
          }

          submitBtn.disabled = true;
          submitBtn.textContent = 'Memproses...';

          try {
            if (isRegister) {
              const res = await apiRequest('/api/v1/auth/register', {
                method: 'POST',
                body: { name: name || 'User', email, password }
              });
              state.token = res.data.token;
              state.user = res.data.user;
              localStorage.setItem('jwt_token', state.token);
              showToast(`Akun berhasil dibuat! Selamat datang, ${state.user.name}`, 'success');
            } else {
              const res = await apiRequest('/api/v1/auth/login', {
                method: 'POST',
                body: { email, password }
              });
              state.token = res.data.token;
              state.user = res.data.user;
              localStorage.setItem('jwt_token', state.token);
              showToast(`Selamat datang kembali, ${state.user.name}!`, 'success');
            }

            closeModal();
            loadDashboard(state.currentLocationId);
          } catch (err) {
            showToast(err.message, 'error');
          } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = isRegister ? 'Daftar Akun Baru' : 'Masuk ke Akun';
          }
        });
      }
    };

    const openModal = () => {
      renderAuthContent();
      modal.classList.add('active');
    };

    const closeModal = () => modal.classList.remove('active');

    profileBtn?.addEventListener('click', openModal);
    avatarImg?.parentElement?.addEventListener('click', openModal);
    closeBtn?.addEventListener('click', closeModal);
    modal.addEventListener('click', (e) => {
      if (e.target === modal) closeModal();
    });
  }

  // 16. Saved Locations Modal
  function openSavedLocationsModal() {
    const modal = document.getElementById('savedLocationsModal');
    const listEl = document.getElementById('savedLocationsList');
    const closeBtn = document.getElementById('savedLocationsCloseBtn');
    if (!modal || !listEl) return;

    modal.classList.add('active');
    closeBtn?.addEventListener('click', () => modal.classList.remove('active'));
    modal.addEventListener('click', (e) => {
      if (e.target === modal) modal.classList.remove('active');
    });

    if (!state.token) {
      listEl.innerHTML = `
        <div style="padding:16px; text-align:center; color:#8c92a4;">
          <p style="margin-bottom:12px;">Silakan login terlebih dahulu untuk menyimpan kota favorit.</p>
          <button class="auth-submit-btn" id="savedLocLoginBtn" style="width:auto; padding:8px 18px;">Login Sekarang</button>
        </div>
      `;
      document.getElementById('savedLocLoginBtn')?.addEventListener('click', () => {
        modal.classList.remove('active');
        document.getElementById('profileBtn')?.click();
      });
      return;
    }

    listEl.innerHTML = '<div style="padding:14px; text-align:center; color:#8c92a4;">Memuat daftar lokasi...</div>';

    apiRequest('/api/v1/locations/saved')
      .then(res => {
        const items = res.data || [];
        if (items.length === 0) {
          listEl.innerHTML = `
            <div style="padding:16px; text-align:center; color:#8c92a4;">
              <p>Belum ada kota favorit yang disimpan.</p>
              <button class="btn-secondary" id="addFirstCityBtn" style="margin-top:10px;">Cari & Tambah Kota</button>
            </div>
          `;
          document.getElementById('addFirstCityBtn')?.addEventListener('click', () => {
            modal.classList.remove('active');
            document.getElementById('addCityBtn')?.click();
          });
          return;
        }

        listEl.innerHTML = items.map(item => {
          const loc = item.location;
          return `
            <div class="saved-loc-item" data-id="${loc.id}">
              <div>
                <div class="saved-loc-title">
                  ${loc.city}
                  ${item.is_default ? '<span class="default-badge">Utama</span>' : ''}
                </div>
                <div class="saved-loc-subtitle">${loc.state ? loc.state + ', ' : ''}${loc.country}</div>
              </div>
              <button class="btn-secondary view-loc-btn" data-id="${loc.id}" style="padding:6px 12px; font-size:11.5px;">Lihat Cuaca</button>
            </div>
          `;
        }).join('');

        listEl.querySelectorAll('.view-loc-btn').forEach(b => {
          b.addEventListener('click', () => {
            const id = b.getAttribute('data-id');
            modal.classList.remove('active');
            loadDashboard(id);
          });
        });
      })
      .catch(err => {
        listEl.innerHTML = `<div style="padding:12px; color:#e53e3e;">${err.message}</div>`;
      });
  }

  // 17. Setup GPS Refresh Button
  function setupGpsButton() {
    const btn = document.getElementById('gpsRefreshBtn');
    if (!btn) return;

    btn.addEventListener('click', async () => {
      btn.classList.add('spin-animation');
      await loadDashboard(state.currentLocationId);
      setTimeout(() => btn.classList.remove('spin-animation'), 600);
      showToast('Data cuaca terkini berhasil diperbarui', 'info');
    });
  }

  // 18. Setup Chart Tabs (Today / Week)
  function setupChartTabs() {
    const tabBtns = document.querySelectorAll('.chart-tabs .tab-btn');
    tabBtns.forEach(btn => {
      btn.addEventListener('click', async () => {
        const period = btn.getAttribute('data-period');
        tabBtns.forEach(b => b.classList.toggle('active', b === btn));
        state.chartPeriod = period;

        if (period === 'week') {
          // Fetch 7-day forecast
          try {
            const res = await apiRequest(`/api/v1/weather/forecast?days=7&unit=${state.unit.toLowerCase()}${state.currentLocationId ? '&location_id=' + state.currentLocationId : ''}`);
            if (res && res.data && res.data.forecasts) {
              const pts = res.data.forecasts.map(f => f.temp_high);
              const lbls = res.data.forecasts.map(f => f.day.substring(0, 3));
              renderSvgChart(pts, lbls, state.unit === 'Fahrenheit' ? '°F' : '°C');
            }
          } catch (e) {
            console.warn('Week forecast error:', e);
          }
        } else {
          // Reload today hourly
          if (state.dashboardData && state.dashboardData.hourly_temperature) {
            const pts = state.dashboardData.hourly_temperature.map(h => h.temp);
            const lbls = state.dashboardData.hourly_temperature.map(h => h.time);
            renderSvgChart(pts, lbls, state.unit === 'Fahrenheit' ? '°F' : '°C');
          } else {
            loadDashboard(state.currentLocationId);
          }
        }
      });
    });
  }

  // 19. Two-City Comparison Logic
  function renderDualSvgChart(series1, series2, labels, name1, name2) {
    const svg = document.getElementById('compareHourlyChartSvg');
    if (!svg || !series1 || !series2 || series1.length === 0 || series2.length === 0) return;

    const width = 600;
    const height = 130;
    const paddingX = 40;
    const paddingTop = 30;
    const paddingBottom = 26;

    const pts1 = series1.map(h => h.temp);
    const pts2 = series2.map(h => h.temp);
    const allPts = [...pts1, ...pts2];

    const min = Math.min(...allPts);
    const max = Math.max(...allPts);
    const range = (max - min) === 0 ? 1 : (max - min);

    const stepX = (width - paddingX * 2) / (series1.length - 1);

    const coords1 = pts1.map((val, idx) => {
      const x = paddingX + idx * stepX;
      const normalized = (val - min) / range;
      const y = (height - paddingBottom) - (normalized * (height - paddingTop - paddingBottom));
      return { x, y, val };
    });

    const coords2 = pts2.map((val, idx) => {
      const x = paddingX + idx * stepX;
      const normalized = (val - min) / range;
      const y = (height - paddingBottom) - (normalized * (height - paddingTop - paddingBottom));
      return { x, y, val };
    });

    function buildSmoothPath(pts) {
      if (pts.length < 2) return '';
      let d = `M ${pts[0].x} ${pts[0].y}`;
      for (let i = 0; i < pts.length - 1; i++) {
        const p0 = pts[i === 0 ? i : i - 1];
        const p1 = pts[i];
        const p2 = pts[i + 1];
        const p3 = pts[i + 2] || p2;

        const cp1x = p1.x + (p2.x - p0.x) / 6;
        const cp1y = p1.y + (p2.y - p0.y) / 6;
        const cp2x = p2.x - (p3.x - p1.x) / 6;
        const cp2y = p2.y - (p3.y - p1.y) / 6;

        d += ` C ${cp1x.toFixed(1)} ${cp1y.toFixed(1)}, ${cp2x.toFixed(1)} ${cp2y.toFixed(1)}, ${p2.x.toFixed(1)} ${p2.y.toFixed(1)}`;
      }
      return d;
    }

    const curve1 = buildSmoothPath(coords1);
    const area1 = `${curve1} L ${coords1[coords1.length - 1].x} ${height} L ${coords1[0].x} ${height} Z`;

    const curve2 = buildSmoothPath(coords2);
    const area2 = `${curve2} L ${coords2[coords2.length - 1].x} ${height} L ${coords2[0].x} ${height} Z`;

    let markers1 = '';
    coords1.forEach((pt) => {
      markers1 += `
        <g class="chart-point-group">
          <circle cx="${pt.x.toFixed(1)}" cy="${pt.y.toFixed(1)}" r="4" fill="#ffffff" stroke="#8c7cf0" stroke-width="2.5"></circle>
          <text x="${pt.x.toFixed(1)}" y="${(pt.y - 8).toFixed(1)}" text-anchor="middle" font-size="10.5" font-weight="700" fill="#6a57e3" font-family="inherit">${Math.round(pt.val)}°</text>
        </g>
      `;
    });

    let markers2 = '';
    coords2.forEach((pt) => {
      markers2 += `
        <g class="chart-point-group">
          <circle cx="${pt.x.toFixed(1)}" cy="${pt.y.toFixed(1)}" r="4" fill="#ffffff" stroke="#00b4d8" stroke-width="2.5"></circle>
          <text x="${pt.x.toFixed(1)}" y="${(pt.y + 16).toFixed(1)}" text-anchor="middle" font-size="10.5" font-weight="700" fill="#028090" font-family="inherit">${Math.round(pt.val)}°</text>
        </g>
      `;
    });

    svg.innerHTML = `
      <defs>
        <linearGradient id="chartGradientC1" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#8c7cf0" stop-opacity="0.22"/>
          <stop offset="100%" stop-color="#8c7cf0" stop-opacity="0.0"/>
        </linearGradient>
        <linearGradient id="chartGradientC2" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#00b4d8" stop-opacity="0.18"/>
          <stop offset="100%" stop-color="#00b4d8" stop-opacity="0.0"/>
        </linearGradient>
      </defs>
      <path d="${area1}" fill="url(#chartGradientC1)"></path>
      <path d="${area2}" fill="url(#chartGradientC2)"></path>
      <path d="${curve1}" fill="none" stroke="#8c7cf0" stroke-width="2.8" stroke-linecap="round" stroke-linejoin="round"></path>
      <path d="${curve2}" fill="none" stroke="#00b4d8" stroke-width="2.8" stroke-linecap="round" stroke-linejoin="round"></path>
      ${markers1}
      ${markers2}
    `;

    const labelsEl = document.getElementById('compareChartTimeLabels');
    if (labelsEl && labels && labels.length > 0) {
      labelsEl.innerHTML = labels.map(l => `<span>${l}</span>`).join('');
    }

    const leg1 = document.getElementById('legendCity1Name');
    const leg2 = document.getElementById('legendCity2Name');
    if (leg1) leg1.textContent = name1;
    if (leg2) leg2.textContent = name2;
  }

  function renderComparison(data) {
    state.comparisonData = data;
    const c1 = data.city1;
    const c2 = data.city2;
    const comp = data.comparison;
    const unitSymbol = state.unit === 'Fahrenheit' ? '°F' : '°C';

    const headlineEl = document.getElementById('summaryHeadline');
    if (headlineEl) headlineEl.textContent = comp.temp_diff_formatted || 'Perbandingan Cuaca';

    const sublineEl = document.getElementById('summarySubline');
    if (sublineEl) sublineEl.textContent = comp.condition_summary || '';

    const badgeTempEl = document.getElementById('badgeTempDiffText');
    if (badgeTempEl) {
      const sign = comp.temp_difference > 0 ? '+' : '';
      badgeTempEl.textContent = `Selisih: ${sign}${comp.temp_difference}${unitSymbol}`;
    }

    const badgeHumEl = document.getElementById('badgeHumidityText');
    if (badgeHumEl) badgeHumEl.textContent = comp.humidity_comparison || 'Kelembapan';

    const badgeRainEl = document.getElementById('badgeRainierText');
    if (badgeRainEl) badgeRainEl.textContent = `Peluang Hujan: ${comp.rainier_city}`;

    const badgeWindEl = document.getElementById('badgeWindierText');
    if (badgeWindEl) badgeWindEl.textContent = `Angin: ${comp.windier_city}`;

    const advisoryEl = document.getElementById('summaryAdvisoryText');
    if (advisoryEl) advisoryEl.textContent = comp.recommendation || 'Kondisi cuaca terpantau normal.';

    const renderCard = (prefix, cityData) => {
      const loc = cityData.location;
      const cur = cityData.current_weather;
      const high = cityData.highlights;

      const headerEl = document.getElementById(`${prefix}HeroHeader`);
      if (headerEl && loc.image_url) {
        headerEl.style.backgroundImage = `url('${loc.image_url}')`;
      }

      const nameEl = document.getElementById(`${prefix}Name`);
      if (nameEl) nameEl.textContent = loc.city;

      const subEl = document.getElementById(`${prefix}Subtitle`);
      if (subEl) subEl.textContent = `${loc.state ? loc.state + ', ' : ''}${loc.country}`;

      const timeEl = document.getElementById(`${prefix}Time`);
      if (timeEl) timeEl.textContent = cur.date_time || 'Hari ini';

      const tempEl = document.getElementById(`${prefix}Temp`);
      if (tempEl) tempEl.textContent = Math.round(cur.temperature);

      const unitEl = document.getElementById(`${prefix}Unit`);
      if (unitEl) unitEl.textContent = unitSymbol;

      const condEl = document.getElementById(`${prefix}Condition`);
      if (condEl) condEl.textContent = cur.condition;

      const humEl = document.getElementById(`${prefix}Humidity`);
      if (humEl) humEl.textContent = high.humidity || '0%';

      const precipEl = document.getElementById(`${prefix}Precip`);
      if (precipEl) precipEl.textContent = high.precipitation || '0%';

      const windEl = document.getElementById(`${prefix}Wind`);
      if (windEl) windEl.textContent = high.wind_speed || '0 km/h';

      const sunEl = document.getElementById(`${prefix}Sun`);
      if (sunEl) sunEl.textContent = `${high.sunrise || '--'} / ${high.sunset || '--'}`;

      const fListEl = document.getElementById(`${prefix}ForecastList`);
      if (fListEl && cityData.three_days_forecast) {
        fListEl.innerHTML = cityData.three_days_forecast.map(f => `
          <div class="mini-forecast-item">
            <span class="mini-forecast-day">${f.day}</span>
            <span class="mini-forecast-cond">${f.condition}</span>
            <span class="mini-forecast-temps">${Math.round(f.temp_high)}° / ${Math.round(f.temp_low)}°</span>
          </div>
        `).join('');
      }
    };

    renderCard('city1', c1);
    renderCard('city2', c2);

    if (c1.hourly_temperature && c2.hourly_temperature) {
      const labels = c1.hourly_temperature.map(h => h.time);
      renderDualSvgChart(c1.hourly_temperature, c2.hourly_temperature, labels, c1.location.city, c2.location.city);
    }

    const h2hList = document.getElementById('h2hList');
    if (h2hList) {
      const parseNum = (str) => {
        const n = parseFloat((str || '').replace(/[^0-9.-]/g, ''));
        return isNaN(n) ? 0 : n;
      };

      const metrics = [
        {
          name: 'Suhu Terkini',
          val1: `${Math.round(c1.current_weather.temperature)}${unitSymbol}`,
          val2: `${Math.round(c2.current_weather.temperature)}${unitSymbol}`,
          num1: Math.max(0, c1.current_weather.temperature),
          num2: Math.max(0, c2.current_weather.temperature)
        },
        {
          name: 'Kelembapan Udara',
          val1: c1.highlights.humidity,
          val2: c2.highlights.humidity,
          num1: parseNum(c1.highlights.humidity),
          num2: parseNum(c2.highlights.humidity)
        },
        {
          name: 'Kecepatan Angin',
          val1: c1.highlights.wind_speed,
          val2: c2.highlights.wind_speed,
          num1: parseNum(c1.highlights.wind_speed),
          num2: parseNum(c2.highlights.wind_speed)
        },
        {
          name: 'Peluang Curah Hujan',
          val1: c1.highlights.precipitation,
          val2: c2.highlights.precipitation,
          num1: parseNum(c1.highlights.precipitation),
          num2: parseNum(c2.highlights.precipitation)
        }
      ];

      h2hList.innerHTML = metrics.map(m => {
        const total = (m.num1 + m.num2) || 1;
        const pct1 = Math.round((m.num1 / total) * 100);
        const pct2 = 100 - pct1;
        return `
          <div class="h2h-item">
            <div class="h2h-meta">
              <span class="h2h-name">${m.name}</span>
              <div class="h2h-values">
                <span class="val-c1">${c1.location.city}: ${m.val1}</span>
                <span style="color:#b2b6cb;">|</span>
                <span class="val-c2">${c2.location.city}: ${m.val2}</span>
              </div>
            </div>
            <div class="h2h-bar-track" title="${c1.location.city}: ${pct1}% vs ${c2.location.city}: ${pct2}%">
              <div class="h2h-fill-c1" style="width: ${pct1}%;"></div>
              <div class="h2h-fill-c2" style="width: ${pct2}%;"></div>
            </div>
          </div>
        `;
      }).join('');
    }
  }

  async function loadComparison(city1 = null, city2 = null) {
    const c1 = city1 || state.compareCity1 || 'Los Angeles';
    const c2 = city2 || state.compareCity2 || 'Jakarta';
    state.compareCity1 = c1;
    state.compareCity2 = c2;

    const select1 = document.getElementById('compareCity1Select');
    const select2 = document.getElementById('compareCity2Select');
    if (select1) select1.value = c1;
    if (select2) select2.value = c2;

    try {
      showToast(`Membandingkan cuaca ${c1} dan ${c2}...`, 'info');
      const unitParam = state.unit.toLowerCase();
      const res = await apiRequest(`/api/v1/weather/compare?city1=${encodeURIComponent(c1)}&city2=${encodeURIComponent(c2)}&unit=${unitParam}`);
      if (res && res.data) {
        renderComparison(res.data);
      }
    } catch (err) {
      showToast('Gagal memuat perbandingan kota: ' + err.message, 'error');
    }
  }

  function showCompareView(city1 = null, city2 = null) {
    state.isCompareMode = true;
    state.isCalendarMode = false;
    const dashGrid = document.getElementById('dashboardGrid');
    const compGrid = document.getElementById('comparisonGrid');
    const calGrid = document.getElementById('calendarGrid');

    if (dashGrid) dashGrid.style.display = 'none';
    if (calGrid) calGrid.style.display = 'none';
    if (compGrid) compGrid.style.display = 'flex';

    document.querySelectorAll('.sidebar-nav .nav-item').forEach(item => {
      item.classList.toggle('active', item.getAttribute('data-tab') === 'compare');
    });

    const c1 = city1 || state.dashboardData?.location?.city || state.compareCity1 || 'Los Angeles';
    const c2 = city2 || (c1 === 'Jakarta' ? 'Tokyo' : 'Jakarta');
    loadComparison(c1, c2);
    compGrid?.scrollIntoView({ behavior: 'smooth' });
  }

  function hideCompareView() {
    state.isCompareMode = false;
    const dashGrid = document.getElementById('dashboardGrid');
    const compGrid = document.getElementById('comparisonGrid');
    const calGrid = document.getElementById('calendarGrid');

    if (compGrid) compGrid.style.display = 'none';
    if (calGrid) calGrid.style.display = 'none';
    if (dashGrid) dashGrid.style.display = 'grid';

    document.querySelectorAll('.sidebar-nav .nav-item').forEach(item => {
      item.classList.toggle('active', item.getAttribute('data-tab') === 'home');
    });
  }

  function setupCompareFeature() {
    const compareNavBtn = document.getElementById('compareNavBtn');
    const headerCompareBtn = document.getElementById('headerCompareBtn');
    const heroCompareBtn = document.getElementById('heroCompareBtn');
    const closeCompareBtn = document.getElementById('closeCompareBtn');
    const swapCitiesBtn = document.getElementById('swapCitiesBtn');
    const applyCompareBtn = document.getElementById('applyCompareBtn');
    const select1 = document.getElementById('compareCity1Select');
    const select2 = document.getElementById('compareCity2Select');

    compareNavBtn?.addEventListener('click', (e) => {
      e.preventDefault();
      showCompareView();
    });

    headerCompareBtn?.addEventListener('click', () => {
      showCompareView();
    });

    heroCompareBtn?.addEventListener('click', () => {
      const currentCity = state.dashboardData?.location?.city || 'Los Angeles';
      showCompareView(currentCity, currentCity === 'Jakarta' ? 'Tokyo' : 'Jakarta');
    });

    closeCompareBtn?.addEventListener('click', () => {
      hideCompareView();
    });

    swapCitiesBtn?.addEventListener('click', () => {
      if (select1 && select2) {
        const temp = select1.value;
        select1.value = select2.value;
        select2.value = temp;
        loadComparison(select1.value, select2.value);
      }
    });

    applyCompareBtn?.addEventListener('click', () => {
      if (select1 && select2) {
        loadComparison(select1.value, select2.value);
      }
    });

    select1?.addEventListener('change', () => {
      if (select1 && select2) loadComparison(select1.value, select2.value);
    });

    select2?.addEventListener('change', () => {
      if (select1 && select2) loadComparison(select1.value, select2.value);
    });

    document.querySelectorAll('#comparisonGrid .preset-pill').forEach(pill => {
      pill.addEventListener('click', () => {
        const c1 = pill.getAttribute('data-c1');
        const c2 = pill.getAttribute('data-c2');
        if (c1 && c2) {
          if (select1) select1.value = c1;
          if (select2) select2.value = c2;
          loadComparison(c1, c2);
        }
      });
    });
  }

  // =========================================================
  // 20. MONTHLY WEATHER CALENDAR FEATURE
  // =========================================================

  const MONTH_NAMES_ID = [
    '', 'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];

  const WEEKDAY_NAMES_ID = {
    'Sunday': 'Minggu',
    'Monday': 'Senin',
    'Tuesday': 'Selasa',
    'Wednesday': 'Rabu',
    'Thursday': 'Kamis',
    'Friday': 'Jumat',
    'Saturday': 'Sabtu'
  };

  function showCalendarView(city = null, year = null, month = null) {
    state.isCalendarMode = true;
    state.isCompareMode = false;

    const dashGrid = document.getElementById('dashboardGrid');
    const compGrid = document.getElementById('comparisonGrid');
    const calGrid = document.getElementById('calendarGrid');

    if (dashGrid) dashGrid.style.display = 'none';
    if (compGrid) compGrid.style.display = 'none';
    if (calGrid) calGrid.style.display = 'flex';

    document.querySelectorAll('.sidebar-nav .nav-item').forEach(item => {
      item.classList.toggle('active', item.getAttribute('data-tab') === 'calendar');
    });

    const targetCity = city || state.calendarCity || state.dashboardData?.location?.city || 'Los Angeles';
    const targetYear = year || state.calendarYear || new Date().getFullYear();
    const targetMonth = month || state.calendarMonth || (new Date().getMonth() + 1);

    loadCalendar(targetCity, targetYear, targetMonth);
    calGrid?.scrollIntoView({ behavior: 'smooth' });
  }

  function hideCalendarView() {
    state.isCalendarMode = false;
    const dashGrid = document.getElementById('dashboardGrid');
    const compGrid = document.getElementById('comparisonGrid');
    const calGrid = document.getElementById('calendarGrid');

    if (calGrid) calGrid.style.display = 'none';
    if (compGrid) compGrid.style.display = 'none';
    if (dashGrid) dashGrid.style.display = 'grid';

    document.querySelectorAll('.sidebar-nav .nav-item').forEach(item => {
      item.classList.toggle('active', item.getAttribute('data-tab') === 'home');
    });
  }

  async function loadCalendar(city = null, year = null, month = null) {
    const c = city || state.calendarCity || state.dashboardData?.location?.city || 'Los Angeles';
    const y = year || state.calendarYear || new Date().getFullYear();
    const m = month || state.calendarMonth || (new Date().getMonth() + 1);

    state.calendarCity = c;
    state.calendarYear = y;
    state.calendarMonth = m;

    const select = document.getElementById('calendarCitySelect');
    if (select) select.value = c;

    const unitParam = state.unit.toLowerCase();
    const unitBadge = document.getElementById('calendarUnitBadge');
    if (unitBadge) {
      unitBadge.textContent = state.unit === 'Fahrenheit' ? 'Fahrenheit (°F)' : 'Celsius (°C)';
    }

    try {
      showToast(`Memuat kalender cuaca ${c} (${MONTH_NAMES_ID[m] || m} ${y})...`, 'info');
      const res = await apiRequest(`/api/v1/weather/calendar?city=${encodeURIComponent(c)}&year=${y}&month=${m}&unit=${unitParam}`);
      if (res && res.data) {
        state.calendarData = res.data;
        renderCalendar(res.data);
      }
    } catch (err) {
      showToast('Gagal memuat kalender: ' + err.message, 'error');
    }
  }

  function renderCalendar(data) {
    if (!data || !data.days) return;

    const monthName = MONTH_NAMES_ID[data.month] || data.month_name || `Bulan ${data.month}`;
    const monthYearTitle = document.getElementById('calendarMonthYearTitle');
    if (monthYearTitle) {
      monthYearTitle.textContent = `${monthName} ${data.year}`;
    }

    const headline = document.getElementById('calendarSummaryHeadline');
    if (headline) {
      headline.textContent = `Prakiraan Cuaca ${data.city} - ${monthName} ${data.year}`;
    }

    const subline = document.getElementById('calendarSummarySubline');
    if (subline) {
      subline.textContent = `Ringkasan analisis cuaca bulanan dan rekomendasi aktivitas`;
    }

    // Monthly Analytics Computation
    const days = data.days;
    let sumHigh = 0;
    let sumLow = 0;
    let sumRain = 0;
    let sunnyCount = 0;
    let rainyCount = 0;

    days.forEach(d => {
      sumHigh += d.temp_high;
      sumLow += d.temp_low;
      sumRain += d.chance_of_rain;

      const c = (d.condition || '').toLowerCase();
      if (c.includes('sun') || c.includes('clear')) {
        sunnyCount++;
      } else if (c.includes('rain') || c.includes('thunder') || c.includes('storm') || c.includes('drizzle')) {
        rainyCount++;
      }
    });

    const avgHigh = (sumHigh / days.length).toFixed(1);
    const avgLow = (sumLow / days.length).toFixed(1);
    const avgRain = Math.round(sumRain / days.length);

    const avgTempText = document.getElementById('calAvgTempText');
    if (avgTempText) avgTempText.textContent = `Rata-rata: ↑${avgHigh}° / ↓${avgLow}°`;

    const rainText = document.getElementById('calRainText');
    if (rainText) rainText.textContent = `Peluang Hujan: ${avgRain}%`;

    const sunnyText = document.getElementById('calSunnyText');
    if (sunnyText) sunnyText.textContent = `${sunnyCount} Hari Cerah`;

    const rainyText = document.getElementById('calRainyText');
    if (rainyText) rainyText.textContent = `${rainyCount} Hari Hujan`;

    const advisoryText = document.getElementById('calAdvisoryText');
    if (advisoryText) {
      if (rainyCount >= 10) {
        advisoryText.textContent = `Bulan ${monthName} di ${data.city} memiliki potensi hujan yang relatif sering (${rainyCount} hari). Disarankan selalu membawa payung dan memantau kondisi harian sebelum kegiatan outdoor.`;
      } else if (sunnyCount >= 15) {
        advisoryText.textContent = `Bulan ${monthName} di ${data.city} didominasi oleh hari-hari cerah (${sunnyCount} hari). Cuaca sangat kondusif untuk bepergian, liburan, dan olahraga di luar ruangan.`;
      } else {
        advisoryText.textContent = `Kondisi cuaca bulan ${monthName} di ${data.city} bervariasi dengan rata-rata suhu harian berkisar ${avgLow}° - ${avgHigh}°.`;
      }
    }

    // Days Grid Rendering
    const gridContainer = document.getElementById('calendarDaysGrid');
    if (!gridContainer) return;
    gridContainer.innerHTML = '';

    // Calculate weekday of the first day (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
    const firstDayDate = new Date(data.year, data.month - 1, 1);
    const startWeekday = firstDayDate.getDay();

    // Render empty padding cells for previous days of the week
    for (let i = 0; i < startWeekday; i++) {
      const pad = document.createElement('div');
      pad.className = 'calendar-day-cell empty-pad';
      gridContainer.appendChild(pad);
    }

    // Find default selected day (today if found in this month, otherwise first day)
    let initialIndex = days.findIndex(d => d.is_today);
    if (initialIndex === -1) {
      initialIndex = 0;
    }
    state.selectedDayIndex = initialIndex;

    days.forEach((day, index) => {
      const cell = document.createElement('div');
      cell.className = `calendar-day-cell ${day.is_today ? 'is-today' : ''} ${index === state.selectedDayIndex ? 'is-selected' : ''}`;
      cell.setAttribute('data-index', index);

      const isTodayPill = day.is_today ? `<span class="day-cell-today-pill">Hari Ini</span>` : '';
      const iconSvg = getWeatherIconSvg(day.condition);

      cell.innerHTML = `
        <div class="day-cell-header">
          <span class="day-cell-number">${day.day}</span>
          ${isTodayPill}
        </div>
        <div class="day-cell-icon">
          ${iconSvg}
        </div>
        <div class="day-cell-cond" title="${day.condition}">${day.condition}</div>
        <div class="day-cell-temp">
          <span class="temp-up">↑${Math.round(day.temp_high)}°</span>
          <span class="temp-down">↓${Math.round(day.temp_low)}°</span>
        </div>
        <div class="day-cell-rain" title="Peluang Hujan: ${day.chance_of_rain}%">💧 ${day.chance_of_rain}%</div>
      `;

      cell.addEventListener('click', () => {
        state.selectedDayIndex = index;
        document.querySelectorAll('.calendar-day-cell').forEach(c => c.classList.remove('is-selected'));
        cell.classList.add('is-selected');
        renderSelectedDayDetail(day, data);
      });

      gridContainer.appendChild(cell);
    });

    // Render inspector for the initially selected day
    if (days[state.selectedDayIndex]) {
      renderSelectedDayDetail(days[state.selectedDayIndex], data);
    }
  }

  function renderSelectedDayDetail(day, data) {
    if (!day) return;

    const monthName = MONTH_NAMES_ID[data.month] || data.month_name || `Bulan ${data.month}`;
    const weekdayId = WEEKDAY_NAMES_ID[day.weekday] || day.weekday;

    const dateTitle = document.getElementById('inspectorDateTitle');
    if (dateTitle) dateTitle.textContent = `${day.day} ${monthName} ${data.year}`;

    const weekdayTitle = document.getElementById('inspectorWeekdayTitle');
    if (weekdayTitle) weekdayTitle.textContent = `${weekdayId}, ${data.city}`;

    const todayPill = document.getElementById('inspectorTodayPill');
    if (todayPill) todayPill.style.display = day.is_today ? 'inline-block' : 'none';

    const iconEl = document.getElementById('inspectorIcon');
    if (iconEl) iconEl.innerHTML = getWeatherIconSvg(day.condition, true);

    const highEl = document.getElementById('inspectorHigh');
    if (highEl) highEl.textContent = `↑${Math.round(day.temp_high)}°`;

    const lowEl = document.getElementById('inspectorLow');
    if (lowEl) lowEl.textContent = `↓${Math.round(day.temp_low)}°`;

    const condEl = document.getElementById('inspectorCondition');
    if (condEl) condEl.textContent = day.condition;

    const rainVal = document.getElementById('inspectorRainVal');
    if (rainVal) rainVal.textContent = `${day.chance_of_rain}%`;

    const rainBar = document.getElementById('inspectorRainBar');
    if (rainBar) rainBar.style.width = `${day.chance_of_rain}%`;

    const advisoryText = document.getElementById('inspectorAdvisoryText');
    if (advisoryText) {
      advisoryText.textContent = day.advisory || 'Kondisi cuaca terpantau normal. Tetap waspada terhadap perubahan cuaca sewaktu-waktu.';
    }
  }

  function setupCalendarFeature() {
    const calendarNavBtn = document.getElementById('calendarNavBtn');
    const headerCalendarBtn = document.getElementById('headerCalendarBtn');
    const heroCalendarBtn = document.getElementById('heroCalendarBtn');
    const closeCalendarBtn = document.getElementById('closeCalendarBtn');
    const prevMonthBtn = document.getElementById('calendarPrevMonthBtn');
    const nextMonthBtn = document.getElementById('calendarNextMonthBtn');
    const currentMonthBtn = document.getElementById('calendarCurrentMonthBtn');
    const citySelect = document.getElementById('calendarCitySelect');

    calendarNavBtn?.addEventListener('click', (e) => {
      e.preventDefault();
      showCalendarView();
    });

    headerCalendarBtn?.addEventListener('click', () => {
      showCalendarView();
    });

    heroCalendarBtn?.addEventListener('click', () => {
      const currentCity = state.dashboardData?.location?.city || 'Los Angeles';
      showCalendarView(currentCity);
    });

    closeCalendarBtn?.addEventListener('click', () => {
      hideCalendarView();
    });

    prevMonthBtn?.addEventListener('click', () => {
      let y = state.calendarYear;
      let m = state.calendarMonth - 1;
      if (m < 1) {
        m = 12;
        y--;
      }
      loadCalendar(state.calendarCity, y, m);
    });

    nextMonthBtn?.addEventListener('click', () => {
      let y = state.calendarYear;
      let m = state.calendarMonth + 1;
      if (m > 12) {
        m = 1;
        y++;
      }
      loadCalendar(state.calendarCity, y, m);
    });

    currentMonthBtn?.addEventListener('click', () => {
      const now = new Date();
      loadCalendar(state.calendarCity, now.getFullYear(), now.getMonth() + 1);
    });

    citySelect?.addEventListener('change', () => {
      loadCalendar(citySelect.value);
    });

    document.querySelectorAll('.calendar-preset-pill').forEach(pill => {
      pill.addEventListener('click', () => {
        const city = pill.getAttribute('data-city');
        if (city) {
          if (citySelect) citySelect.value = city;
          loadCalendar(city);
        }
      });
    });
  }

  // 21. Setup Sidebar Navigation Links
  function setupSidebarNav() {
    const navItems = document.querySelectorAll('.sidebar-nav .nav-item');
    navItems.forEach(item => {
      item.addEventListener('click', (e) => {
        const tab = item.getAttribute('data-tab');
        if (tab === 'settings') return; // Handled by settings modal
        e.preventDefault();

        navItems.forEach(i => i.classList.remove('active'));
        item.classList.add('active');

        if (tab === 'calendar') {
          showCalendarView();
        } else if (tab === 'compare') {
          showCompareView();
        } else {
          if (state.isCompareMode) {
            hideCompareView();
          }
          if (state.isCalendarMode) {
            hideCalendarView();
          }
          if (tab === 'locations') {
            openSavedLocationsModal();
          } else if (tab === 'home') {
            state.currentLocationId = null;
            loadDashboard();
          } else if (tab === 'forecast') {
            document.querySelector('.forecast-section')?.scrollIntoView({ behavior: 'smooth' });
          } else if (tab === 'analytics') {
            document.querySelector('.hourly-chart-section')?.scrollIntoView({ behavior: 'smooth' });
          }
        }
      });
    });
  }

  // 20. Initialize User Session if Token Exists
  async function initAuth() {
    if (!state.token) return;

    try {
      const res = await apiRequest('/api/v1/users/profile');
      if (res && res.data) {
        state.user = res.data;
        if (state.user.temperature_unit) {
          state.unit = state.user.temperature_unit;
          localStorage.setItem('temperature_unit', state.unit);
        }
      }
    } catch (err) {
      console.warn('Session expired or invalid:', err.message);
      localStorage.removeItem('jwt_token');
      state.token = null;
      state.user = null;
    }
  }

  // 21. Application Bootstrap
  async function init() {
    console.log('[go_vibe] Initializing Weather & Analytics Dashboard...');
    await initAuth();
    setupSearchBar();
    setupCityModal();
    setupSettingsModal();
    setupAuthModal();
    setupGpsButton();
    setupChartTabs();
    setupSidebarNav();
    setupCompareFeature();
    setupCalendarFeature();

    // Initial Dashboard Load
    await loadDashboard();
    console.log('[go_vibe] Dashboard successfully connected to backend!');
  }

  // Run on DOM Ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();

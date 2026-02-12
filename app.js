// Check authentication on page load
async function checkAuth() {
  const token = localStorage.getItem('token');
  const currentPath = window.location.pathname;
  const publicPaths = ['/', '/auth', '/profile'];
  
  // Allow access to public paths
  if (publicPaths.includes(currentPath)) {
    return;
  }
  
  // Check if user is authenticated
  if (!token) {
    window.location = '/auth';
    return;
  }
  
  // Verify token is still valid
  try {
    await fetchJSON('/api/profile/me', {
      headers: { 'Authorization': 'Bearer ' + token }
    });
  } catch (err) {
    // Token invalid, redirect to auth
    localStorage.removeItem('token');
    window.location = '/auth';
  }
}

// Apply favorite team theme on every page load
window.addEventListener('DOMContentLoaded', async () => {
  // Check authentication first
  await checkAuth();
  
  const saved = localStorage.getItem('favTheme');
  if (saved) {
    document.body.className = 'theme-' + saved;
  }
  
  // Set active navigation link based on current page
  const currentPath = window.location.pathname;
  document.querySelectorAll('#sidebar a').forEach(link => {
    if (link.getAttribute('href') === currentPath || 
        (currentPath === '/' && link.getAttribute('href') === '/feed')) {
      link.classList.add('active');
    } else {
      link.classList.remove('active');
    }
  });
  
  // Add smooth scroll behavior
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute('href'));
      if (target) {
        target.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
    });
  });
  
  // Load user profile and update header
  const token = localStorage.getItem('token');
  if (token && document.querySelector('.user-profile')) {
    try {
      const me = await fetchJSON('/api/profile/me', {
        headers: { 'Authorization': 'Bearer ' + token }
      });
      const avatar = document.querySelector('.user-avatar');
      const nameSpan = document.querySelector('.user-profile span');
      if (avatar && me.Name) {
        const initials = me.Name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2);
        avatar.textContent = initials;
      }
      if (nameSpan && me.Name) {
        nameSpan.textContent = me.Name;
      }
    } catch (err) {
      console.error('Failed to load user profile:', err);
    }
  }
});

function getToken() {
  return localStorage.getItem('token') || '';
}

async function logout() {
  try {
    // Inform backend to clear HTTP-only auth cookie
    await fetch('/api/auth/logout', { method: 'POST' });
  } catch (e) {
    console.error('Logout request failed:', e);
  }
  localStorage.removeItem('token');
  localStorage.removeItem('favTheme');
  window.location = '/auth';
}

function setThemeFromTeam(team) {
  if (!team) return;
  document.documentElement.style.setProperty('--primary', team.primaryColor || '#222');
  document.documentElement.style.setProperty('--secondary', team.secondaryColor || '#555');
}

async function fetchJSON(url, opts={}) {
  const res = await fetch(url, opts);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

// Utility function to format numbers
function formatNumber(num) {
  return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

// Utility function to show notifications
function showNotification(message, type = 'info') {
  const notification = document.createElement('div');
  notification.style.cssText = `
    position: fixed;
    top: 90px;
    right: 20px;
    padding: 16px 24px;
    background: var(--bg-secondary);
    border-left: 4px solid var(--accent-primary);
    border-radius: 8px;
    box-shadow: var(--shadow-lg);
    z-index: 1000;
    animation: slideIn 0.3s ease;
  `;
  notification.textContent = message;
  document.body.appendChild(notification);
  
  setTimeout(() => {
    notification.style.animation = 'slideOut 0.3s ease';
    setTimeout(() => notification.remove(), 300);
  }, 3000);
}

// Add CSS animations for notifications
const style = document.createElement('style');
style.textContent = `
  @keyframes slideIn {
    from { transform: translateX(400px); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
  @keyframes slideOut {
    from { transform: translateX(0); opacity: 1; }
    to { transform: translateX(400px); opacity: 0; }
  }
`;
document.head.appendChild(style);

window.initTableLive = function initTableLive() {
  const tbody = document.getElementById('table-body');
  const liveIndicator = document.getElementById('live-indicator');
  let ws;
  function renderTable(rows) {
    tbody.innerHTML = '';
    rows.forEach((r, idx) => {
      const tr = document.createElement('tr');
      let statusCell = '';
      if (r.live) {
        tr.classList.add('team-live');
        statusCell = '<span class="live-badge">LIVE</span>';
      }
      tr.innerHTML = `<td>${idx+1}</td><td>${r.team}</td><td>${r.played}</td><td>${r.points}</td><td>${r.gd}</td><td>${statusCell}</td>`;
      tbody.appendChild(tr);
    });
  }
  function connectWS() {
    ws = new WebSocket('ws://' + window.location.host + '/ws/standings');
    ws.onopen = () => {
      console.log('WebSocket connected');
    };
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.standings) {
        renderTable(data.standings);
        const anyLive = data.standings.some(t => t.live);
        liveIndicator.style.display = anyLive ? '' : 'none';
      }
    };
    ws.onclose = () => {
      setTimeout(connectWS, 2000);
    };
    ws.onerror = () => {
      ws.close();
    };
  }
  connectWS();
}

window.initProfilePage = async function initProfilePage() {
  const loginBtn = document.getElementById('loginBtn');
  const emailEl = document.getElementById('email');
  const pwEl = document.getElementById('password');
  const statusEl = document.getElementById('loginStatus');
  const teamSelect = document.getElementById('teamSelect');
  const saveFavBtn = document.getElementById('saveFavBtn');
  const meData = document.getElementById('meData');
  const nameDisplay = document.getElementById('displayName');
  const emailDisplay = document.getElementById('displayEmail');
  const roleDisplay = document.getElementById('displayRole');
  const favTeamDisplay = document.getElementById('displayFavoriteTeam');

  loginBtn.onclick = async () => {
    try {
      const body = { email: emailEl.value, password: pwEl.value };
      const resp = await fetchJSON('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      });
      localStorage.setItem('token', resp.token);
      statusEl.textContent = 'Logged in';
      await loadMe();
    } catch (e) {
      statusEl.textContent = 'Login failed';
    }
  };

  async function loadTeams() {
    const teams = await fetchJSON('/api/teams');
    teamSelect.innerHTML = '<option value="">Select a team...</option>';
    teams.forEach(t => {
      const opt = document.createElement('option');
      opt.value = t.ID || t.id || t.Id;
      opt.textContent = t.Name || t.name;
      opt.dataset.primaryColor = t.PrimaryColor || t.primaryColor;
      opt.dataset.secondaryColor = t.SecondaryColor || t.secondaryColor;
      teamSelect.appendChild(opt);
    });
  }

  async function loadMe() {
    try {
      const me = await fetchJSON('/api/profile/me', {
        headers: { 'Authorization': 'Bearer ' + getToken() }
      });
      if (meData) {
        meData.dataset.raw = JSON.stringify(me, null, 2);
      }

      // Update profile summary fields
      if (nameDisplay && (me.Name || me.name)) {
        nameDisplay.textContent = me.Name || me.name;
      }
      if (emailDisplay && (me.Email || me.email)) {
        emailDisplay.textContent = me.Email || me.email;
      }
      if (roleDisplay && (me.Role || me.role)) {
        roleDisplay.textContent = me.Role || me.role;
      }
      const favTeam = me.FavoriteTeam || me.favoriteTeam;
      if (favTeamDisplay) {
        favTeamDisplay.textContent = favTeam && (favTeam.Name || favTeam.name) ? (favTeam.Name || favTeam.name) : 'Not set';
      }

      // Pre-select current favorite team in dropdown if available
      if (teamSelect && favTeam && (favTeam.ID || favTeam.Id || favTeam.id)) {
        const favId = String(favTeam.ID || favTeam.Id || favTeam.id);
        const option = Array.from(teamSelect.options).find(o => o.value === favId);
        if (option) {
          teamSelect.value = favId;
        }
      }

      // Apply visual theme based on favorite team
      setThemeFromTeam(favTeam);
      if (favTeam && (favTeam.Name || favTeam.name)) {
        const name = (favTeam.Name || favTeam.name).toLowerCase();
        let themeKey = '';
        if (name.includes('manchester united')) themeKey = 'manunited';
        else if (name.includes('manchester city')) themeKey = 'mancity';
        else if (name.includes('liverpool')) themeKey = 'liverpool';
        else if (name.includes('arsenal')) themeKey = 'arsenal';
        else if (name.includes('tottenham')) themeKey = 'tottenham';
        else if (name.includes('chelsea')) themeKey = 'chelsea';
        else if (name.includes('newcastle')) themeKey = 'newcastle';
        else if (name.includes('aston villa')) themeKey = 'astonvilla';

        if (themeKey) {
          localStorage.setItem('favTheme', themeKey);
          document.body.className = 'theme-' + themeKey;
        }
      }
    } catch (e) {
      if (meData) {
        meData.dataset.raw = 'Not logged in';
      }
    }
  }

  saveFavBtn.onclick = async () => {
    const teamId = parseInt(teamSelect.value, 10);
    if (!teamId) {
      showNotification('Please select a team first.', 'info');
      return;
    }
    try {
      await fetchJSON('/api/profile/favorite', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer ' + getToken()
        },
        body: JSON.stringify({ teamId }),
      });
      await loadMe();
      showNotification('Favorite team saved successfully.', 'info');
    } catch (e) {
      showNotification('Failed to save favorite team.', 'error');
    }
  };

  await loadTeams();
  await loadMe();
}

window.initAuthPage = async function initAuthPage() {
  const tabs = document.querySelectorAll('.auth-tab');
  const loginForm = document.getElementById('loginForm');
  const registerForm = document.getElementById('registerForm');
  const loginMsg = document.getElementById('loginMessage');
  const registerMsg = document.getElementById('registerMessage');
  const teamSelect = document.getElementById('registerTeam');

  // Load teams for registration
  async function loadTeams() {
    try {
      const teams = await fetchJSON('/api/teams');
      teamSelect.innerHTML = '<option value="">Select a team...</option>';
      teams.forEach(t => {
        const opt = document.createElement('option');
        opt.value = t.ID || t.id || t.Id;
        opt.textContent = t.Name || t.name;
        teamSelect.appendChild(opt);
      });
    } catch (err) {
      console.error('Failed to load teams:', err);
    }
  }
  await loadTeams();

  // Tab switching
  tabs.forEach(tab => {
    tab.addEventListener('click', () => {
      tabs.forEach(t => t.classList.remove('auth-tab--active'));
      tab.classList.add('auth-tab--active');
      const target = tab.dataset.tab;
      if (target === 'login') {
        loginForm.classList.remove('auth-form--hidden');
        registerForm.classList.add('auth-form--hidden');
      } else {
        registerForm.classList.remove('auth-form--hidden');
        loginForm.classList.add('auth-form--hidden');
      }
    });
  });

  function setMessage(el, msg, type) {
    el.textContent = msg;
    el.classList.remove('auth-message--success', 'auth-message--error');
    if (type === 'success') el.classList.add('auth-message--success');
    if (type === 'error') el.classList.add('auth-message--error');
  }

  async function checkUserAndRedirect(token) {
    try {
      const me = await fetchJSON('/api/profile/me', {
        headers: { 'Authorization': 'Bearer ' + token }
      });
      // Check if user has favorite team
      if (!me.FavoriteTeam && !me.favoriteTeam) {
        window.location = '/profile';
      } else {
        window.location = '/feed';
      }
    } catch (err) {
      console.error('Failed to check user:', err);
      window.location = '/profile';
    }
  }

  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const email = document.getElementById('loginEmail').value.trim();
      const password = document.getElementById('loginPassword').value;
      try {
        const resp = await fetchJSON('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ email, password }),
        });
        localStorage.setItem('token', resp.token);
        setMessage(loginMsg, 'Login successful! Redirecting…', 'success');
        showNotification('Welcome back!', 'info');
        setTimeout(() => {
          checkUserAndRedirect(resp.token);
        }, 800);
      } catch (err) {
        console.error(err);
        setMessage(loginMsg, 'Login failed. Check your email and password.', 'error');
      }
    });
  }

  if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      const name = document.getElementById('registerName').value.trim();
      const email = document.getElementById('registerEmail').value.trim();
      const password = document.getElementById('registerPassword').value;
      const teamId = parseInt(teamSelect.value, 10);
      
      if (!teamId) {
        setMessage(registerMsg, 'Please select your favorite team.', 'error');
        return;
      }

      try {
        const body = { name, email, password };
        if (teamId) {
          body.favoriteTeam = teamId;
        }
        await fetchJSON('/api/auth/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        setMessage(registerMsg, 'Account created! Logging you in...', 'success');
        showNotification('Account created successfully!', 'info');
        
        // Auto-login after registration
        setTimeout(async () => {
          try {
            const resp = await fetchJSON('/api/auth/login', {
              method: 'POST',
              headers: { 'Content-Type': 'application/json' },
              body: JSON.stringify({ email, password }),
            });
            localStorage.setItem('token', resp.token);
            window.location = '/feed';
          } catch (err) {
            // Switch to login tab if auto-login fails
            tabs.forEach(t => {
              if (t.dataset.tab === 'login') t.classList.add('auth-tab--active');
              else t.classList.remove('auth-tab--active');
            });
            registerForm.classList.add('auth-form--hidden');
            loginForm.classList.remove('auth-form--hidden');
            setMessage(registerMsg, 'Account created! Please log in.', 'success');
          }
        }, 1000);
      } catch (err) {
        console.error(err);
        setMessage(registerMsg, 'Registration failed. Try a different email.', 'error');
      }
    });
  }
}

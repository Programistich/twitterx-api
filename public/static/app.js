// Home page JavaScript

const RECENT_SEARCHES_KEY = 'twitterx_recent_searches';
const MAX_RECENT_SEARCHES = 5;

document.addEventListener('DOMContentLoaded', () => {
    const searchForm = document.getElementById('searchForm');
    const usernameInput = document.getElementById('usernameInput');
    const searchBtn = document.getElementById('searchBtn');
    const recentSearches = document.getElementById('recentSearches');
    const recentList = document.getElementById('recentList');
    const errorModal = new bootstrap.Modal(document.getElementById('errorModal'));
    const errorModalText = document.getElementById('errorModalText');
    const errorModalLabel = document.getElementById('errorModalLabel');

    // Load recent searches
    loadRecentSearches();

    // Handle form submit
    searchForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = usernameInput.value.trim().replace(/^@/, '');

        if (!username) {
            showError('Invalid Username', 'Please enter a username');
            return;
        }

        if (!isValidUsername(username)) {
            showError('Invalid Username', 'Username must be 1-15 characters, letters, numbers and underscores only');
            return;
        }

        // Disable button and show loading
        searchBtn.disabled = true;
        searchBtn.innerHTML = '<span class="spinner-border spinner-border-sm me-2"></span>Checking...';

        try {
            // Check if user exists via API
            const response = await fetch(`/api/users/${encodeURIComponent(username)}`);

            if (!response.ok) {
                if (response.status === 404) {
                    showError('User Not Found', `The user @${username} does not exist or is unavailable.`);
                } else {
                    const errorText = await response.text();
                    showError('Error', errorText || 'Failed to fetch user data');
                }
                return;
            }

            // User exists - save to recent and navigate
            saveRecentSearch(username);
            window.location.href = '/' + encodeURIComponent(username);

        } catch (error) {
            showError('Connection Error', 'Unable to connect to the server. Please try again.');
        } finally {
            // Re-enable button
            searchBtn.disabled = false;
            searchBtn.innerHTML = 'Search';
        }
    });

    function isValidUsername(username) {
        // Twitter usernames: 1-15 chars, alphanumeric and underscores
        return /^[a-zA-Z0-9_]{1,15}$/.test(username);
    }

    function showError(title, message) {
        errorModalLabel.textContent = title;
        errorModalText.textContent = message;
        errorModal.show();
    }

    function loadRecentSearches() {
        const searches = getRecentSearches();

        if (searches.length === 0) {
            recentSearches.classList.add('d-none');
            return;
        }

        recentSearches.classList.remove('d-none');
        recentList.innerHTML = '';

        searches.forEach(username => {
            const item = document.createElement('a');
            item.href = '/' + encodeURIComponent(username);
            item.className = 'btn btn-outline-secondary btn-sm';
            item.textContent = '@' + username;
            recentList.appendChild(item);
        });
    }

    function getRecentSearches() {
        try {
            const data = localStorage.getItem(RECENT_SEARCHES_KEY);
            return data ? JSON.parse(data) : [];
        } catch {
            return [];
        }
    }

    function saveRecentSearch(username) {
        let searches = getRecentSearches();

        // Remove if already exists
        searches = searches.filter(s => s.toLowerCase() !== username.toLowerCase());

        // Add to beginning
        searches.unshift(username);

        // Keep only max items
        searches = searches.slice(0, MAX_RECENT_SEARCHES);

        try {
            localStorage.setItem(RECENT_SEARCHES_KEY, JSON.stringify(searches));
        } catch {
            // Ignore storage errors
        }
    }
});

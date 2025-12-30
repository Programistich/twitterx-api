// Profile page JavaScript

document.addEventListener('DOMContentLoaded', async () => {
    const loadingState = document.getElementById('loadingState');
    const errorState = document.getElementById('errorState');
    const profileContent = document.getElementById('profileContent');
    const errorText = document.getElementById('errorText');

    // Get username from URL path
    const path = window.location.pathname;
    const username = decodeURIComponent(path.substring(1));

    if (!username) {
        window.location.href = '/';
        return;
    }

    try {
        const response = await fetch(`/api/users/${encodeURIComponent(username)}`);

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(response.status === 404 ? 'User not found' : errorData || 'Failed to load user');
        }

        const data = await response.json();

        if (data.code !== 200 || !data.user) {
            throw new Error(data.message || 'User not found');
        }

        displayProfile(data.user);
        loadTweets(username);
    } catch (error) {
        showError(error.message);
    }

    function displayProfile(user) {
        loadingState.classList.add('d-none');
        profileContent.classList.remove('d-none');

        // Page title
        document.title = `${user.name} (@${user.screen_name}) - TwitterX`;

        // Header
        document.getElementById('headerName').textContent = user.name;
        document.getElementById('headerTweets').textContent = formatNumber(user.tweets) + ' posts';

        // Banner
        if (user.banner_url) {
            const banner = document.getElementById('banner');
            banner.src = user.banner_url;
            banner.classList.remove('d-none');
        }

        // Avatar
        const avatar = document.getElementById('avatar');
        avatar.src = user.avatar_url || '';
        avatar.alt = `${user.name}'s avatar`;

        // Twitter link
        document.getElementById('twitterLink').href = user.url || `https://twitter.com/${user.screen_name}`;

        // Name
        document.getElementById('displayName').textContent = user.name;

        // Verification
        if (user.verification) {
            const badge = document.getElementById('verificationBadge');
            badge.classList.remove('d-none');
            if (user.verification.type === 'Business' || user.verification.type === 'Government') {
                badge.querySelector('.verified-icon').classList.add('gold');
            }
        }

        // Protected
        if (user.protected) {
            document.getElementById('protectedBadge').classList.remove('d-none');
        }

        // Username
        document.getElementById('username').textContent = '@' + user.screen_name;

        // Bio
        const bio = document.getElementById('bio');
        if (user.description) {
            bio.textContent = user.description;
        } else {
            bio.classList.add('d-none');
        }

        // Location
        if (user.location) {
            document.getElementById('locationWrapper').classList.remove('d-none');
            document.getElementById('location').textContent = user.location;
        }

        // Website
        if (user.website) {
            document.getElementById('websiteWrapper').classList.remove('d-none');
            const websiteEl = document.getElementById('website');
            websiteEl.href = user.website;
            websiteEl.textContent = formatUrl(user.website);
        }

        // Joined
        if (user.joined) {
            document.getElementById('joinedWrapper').classList.remove('d-none');
            document.getElementById('joined').textContent = 'Joined ' + formatDate(user.joined);
        }

        // Stats
        document.getElementById('following').textContent = formatNumber(user.following);
        document.getElementById('followers').textContent = formatNumber(user.followers);
        document.getElementById('tweets').textContent = formatNumber(user.tweets);
        document.getElementById('likes').textContent = formatNumber(user.likes);
        document.getElementById('media').textContent = formatNumber(user.media_count);
    }

    function showError(message) {
        loadingState.classList.add('d-none');
        errorState.classList.remove('d-none');
        errorText.textContent = message || 'This account doesn\'t exist.';
    }

    function formatNumber(num) {
        if (num === undefined || num === null) return '0';
        if (num >= 1000000) return (num / 1000000).toFixed(1).replace(/\.0$/, '') + 'M';
        if (num >= 1000) return (num / 1000).toFixed(1).replace(/\.0$/, '') + 'K';
        return num.toLocaleString();
    }

    function formatUrl(url) {
        try {
            const parsed = new URL(url);
            let display = parsed.hostname.replace(/^www\./, '');
            if (parsed.pathname && parsed.pathname !== '/') display += parsed.pathname;
            return display.length > 30 ? display.substring(0, 30) + '...' : display;
        } catch {
            return url;
        }
    }

    function formatDate(dateStr) {
        try {
            return new Date(dateStr).toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
        } catch {
            return dateStr;
        }
    }

    async function loadTweets(username) {
        const tweetsLoading = document.getElementById('tweetsLoading');
        const tweetsError = document.getElementById('tweetsError');
        const tweetsErrorText = document.getElementById('tweetsErrorText');
        const tweetsList = document.getElementById('tweetsList');

        try {
            // Fetch tweet IDs
            const response = await fetch(`/api/users/${encodeURIComponent(username)}/tweets`);
            if (!response.ok) {
                throw new Error('Failed to load tweets');
            }

            const data = await response.json();
            const tweetIDs = data.tweet_ids || [];

            if (tweetIDs.length === 0) {
                tweetsLoading.classList.add('d-none');
                tweetsError.classList.remove('d-none');
                tweetsErrorText.textContent = 'No posts found';
                return;
            }

            // Fetch each tweet's details
            const tweets = await Promise.all(
                tweetIDs.map(async (id) => {
                    try {
                        const tweetResponse = await fetch(`/api/users/${encodeURIComponent(username)}/tweets/${id}`);
                        if (!tweetResponse.ok) return null;
                        const tweetData = await tweetResponse.json();
                        return tweetData.tweet || null;
                    } catch {
                        return null;
                    }
                })
            );

            // Filter out failed requests
            const validTweets = tweets.filter(t => t !== null);

            if (validTweets.length === 0) {
                tweetsLoading.classList.add('d-none');
                tweetsError.classList.remove('d-none');
                tweetsErrorText.textContent = 'Failed to load post details';
                return;
            }

            // Display tweets
            tweetsLoading.classList.add('d-none');
            tweetsList.classList.remove('d-none');
            tweetsList.innerHTML = validTweets.map(tweet => createTweetCard(tweet)).join('');

        } catch (error) {
            tweetsLoading.classList.add('d-none');
            tweetsError.classList.remove('d-none');
            tweetsErrorText.textContent = error.message || 'Failed to load posts';
        }
    }

    function createTweetCard(tweet) {
        const date = formatTweetDate(tweet.created_at);
        const mediaHTML = tweet.media ? createMediaHTML(tweet.media) : '';
        const quoteHTML = tweet.quote ? createQuoteHTML(tweet.quote) : '';

        return `
            <div class="card bg-dark border-secondary mb-3">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-2">
                        <small class="text-secondary">${date}</small>
                        <a href="${tweet.url}" target="_blank" class="text-secondary text-decoration-none">
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M14 3v2h3.59l-9.83 9.83 1.41 1.41L19 6.41V10h2V3h-7zm-2 16H5V5h7V3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2v7h-7z"/>
                            </svg>
                        </a>
                    </div>
                    <p class="card-text mb-2" style="white-space: pre-wrap;">${escapeHtml(tweet.text)}</p>
                    ${mediaHTML}
                    ${quoteHTML}
                    <div class="d-flex gap-4 text-secondary small mt-3">
                        <span title="Replies">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" class="me-1">
                                <path d="M21.99 4c0-1.1-.89-2-1.99-2H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h14l4 4-.01-18z"/>
                            </svg>
                            ${formatNumber(tweet.replies)}
                        </span>
                        <span title="Retweets">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" class="me-1">
                                <path d="M7 7h10v3l4-4-4-4v3H5v6h2V7zm10 10H7v-3l-4 4 4 4v-3h12v-6h-2v4z"/>
                            </svg>
                            ${formatNumber(tweet.retweets)}
                        </span>
                        <span title="Likes">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" class="me-1">
                                <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                            </svg>
                            ${formatNumber(tweet.likes)}
                        </span>
                        ${tweet.views ? `
                        <span title="Views">
                            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" class="me-1">
                                <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
                            </svg>
                            ${formatNumber(tweet.views)}
                        </span>
                        ` : ''}
                    </div>
                </div>
            </div>
        `;
    }

    function createMediaHTML(media) {
        if (!media) return '';

        let html = '<div class="mt-2">';

        // Photos
        if (media.photos && media.photos.length > 0) {
            const gridClass = media.photos.length > 1 ? 'd-grid gap-2' : '';
            const gridStyle = media.photos.length > 1 ? 'grid-template-columns: repeat(2, 1fr);' : '';
            html += `<div class="${gridClass}" style="${gridStyle}">`;
            media.photos.forEach(photo => {
                html += `<img src="${photo.url}" alt="Tweet image" class="img-fluid rounded" style="max-height: 300px; object-fit: cover; width: 100%;">`;
            });
            html += '</div>';
        }

        // Videos
        if (media.videos && media.videos.length > 0) {
            media.videos.forEach(video => {
                html += `
                    <video controls class="w-100 rounded" style="max-height: 400px;" poster="${video.thumbnail_url}">
                        <source src="${video.url}" type="video/mp4">
                    </video>
                `;
            });
        }

        html += '</div>';
        return html;
    }

    function createQuoteHTML(quote) {
        if (!quote) return '';

        return `
            <div class="border border-secondary rounded p-3 mt-2">
                <div class="d-flex align-items-center gap-2 mb-2">
                    <img src="${quote.author.avatar_url}" alt="${quote.author.name}" class="rounded-circle" width="20" height="20">
                    <span class="fw-bold">${escapeHtml(quote.author.name)}</span>
                    <span class="text-secondary">@${quote.author.screen_name}</span>
                </div>
                <p class="mb-0 small" style="white-space: pre-wrap;">${escapeHtml(quote.text)}</p>
            </div>
        `;
    }

    function formatTweetDate(dateStr) {
        try {
            const date = new Date(dateStr);
            const now = new Date();
            const diff = now - date;

            // Less than 24 hours
            if (diff < 86400000) {
                const hours = Math.floor(diff / 3600000);
                if (hours < 1) {
                    const mins = Math.floor(diff / 60000);
                    return mins < 1 ? 'now' : `${mins}m`;
                }
                return `${hours}h`;
            }

            // Less than 7 days
            if (diff < 604800000) {
                const days = Math.floor(diff / 86400000);
                return `${days}d`;
            }

            // Same year
            if (date.getFullYear() === now.getFullYear()) {
                return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
            }

            // Different year
            return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
        } catch {
            return dateStr;
        }
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
});

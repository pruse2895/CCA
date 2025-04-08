const API_BASE_URL = 'http://localhost:8080'; // Your Go backend URL

// --- DOM Element References ---
const authSection = document.getElementById('auth-section');
const loggedInSection = document.getElementById('logged-in-section');

// Forms
const loginForm = document.getElementById('login-form');
const signupForm = document.getElementById('signup-form');
const updateProfileForm = document.getElementById('update-profile-form');
const createAnnouncementForm = document.getElementById('create-announcement-form');

// Message Areas
const loginMessage = document.getElementById('login-message');
const signupMessage = document.getElementById('signup-message');
const profileMessage = document.getElementById('profile-message');
const announcementsMessage = document.getElementById('announcements-message');
const createAnnouncementMessage = document.getElementById('create-announcement-message');
const cricketersMessage = document.getElementById('cricketers-message');

// Display Areas
const announcementsList = document.getElementById('announcements-list');
const profileDetails = document.getElementById('profile-details');
const cricketersList = document.getElementById('cricketers-list');

// Buttons
const logoutButton = document.getElementById('logout-button');

// Specific Sections for Visibility Control
const cricketerProfileSection = document.getElementById('cricketer-profile-section');
const adminSection = document.getElementById('admin-section');

// --- State Variables ---
let authToken = localStorage.getItem('authToken'); // Store token locally
let userRole = localStorage.getItem('userRole');
let userId = localStorage.getItem('userId');

// --- Initialization ---
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    checkLoginState(); // Check if already logged in on page load
});

function setupEventListeners() {
    loginForm.addEventListener('submit', handleLogin);
    signupForm.addEventListener('submit', handleSignup);
    logoutButton.addEventListener('click', handleLogout);
    updateProfileForm.addEventListener('submit', handleUpdateProfile);
    createAnnouncementForm.addEventListener('submit', handleCreateAnnouncement);
    // Add more listeners as needed
}

function checkLoginState() {
    if (authToken) {
        // User is logged in, show relevant sections
        authSection.style.display = 'none';
        loggedInSection.style.display = 'block';
        fetchAnnouncements(); // Fetch announcements on load if logged in

        if (userRole === 'cricketer') {
            cricketerProfileSection.style.display = 'block';
            adminSection.style.display = 'none';
            fetchCricketerProfile(); // Fetch profile if cricketer
        } else if (userRole === 'admin') {
            cricketerProfileSection.style.display = 'none';
            adminSection.style.display = 'block';
            fetchAllCricketers(); // Fetch all cricketers if admin
        }
    } else {
        // User is not logged in
        authSection.style.display = 'block';
        loggedInSection.style.display = 'none';
    }
}

// --- Helper Functions ---

// Function to display messages
function showMessage(element, message, isError = false) {
    element.textContent = message;
    element.className = isError ? 'error-message' : '';
}

// Function to make API requests
async function apiRequest(endpoint, method = 'GET', body = null, requiresAuth = false) {
    const url = `${API_BASE_URL}${endpoint}`;
    const headers = new Headers({
        'Content-Type': 'application/json',
    });

    if (requiresAuth && authToken) {
        headers.append('Authorization', `Bearer ${authToken}`);
    }

    const config = {
        method: method,
        headers: headers,
    };

    if (body) {
        config.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(url, config);

        // Handle cases where the response might not have a body (e.g., 201 Created, 204 No Content)
        const contentType = response.headers.get("content-type");
        let data = null;
        if (contentType && contentType.indexOf("application/json") !== -1) {
            data = await response.json();
        } else {
            // If not JSON, maybe just get text or handle based on status
            // For simplicity, we'll just check the status
        }

        if (!response.ok) {
            // Try to get error message from backend response body
            const errorMessage = data?.error || data?.message || `HTTP error! status: ${response.status}`;
            throw new Error(errorMessage);
        }
        return data; // Return parsed JSON data or null
    } catch (error) {
        console.error('API Request Error:', error);
        throw error; // Re-throw the error to be caught by the caller
    }
}

// Placeholder functions for API calls - Will implement next
async function handleLogin(event) {
    event.preventDefault();
    showMessage(loginMessage, ''); // Clear previous messages

    const identifierInput = document.getElementById('login-identifier');
    const passwordInput = document.getElementById('login-password');
    const submitter = event.submitter; // Get the button that triggered the submit
    const loginType = submitter ? submitter.value : null; // 'cricketer' or 'admin'

    if (!loginType) {
        showMessage(loginMessage, 'Could not determine login type.', true);
        return;
    }

    const identifier = identifierInput.value.trim();
    const password = passwordInput.value.trim();

    let endpoint = '';
    let requestBody = {};

    if (loginType === 'cricketer') {
        endpoint = '/api/login';
        requestBody = {
            mobile: identifier,
            password: password
        };
    } else if (loginType === 'admin') {
        endpoint = '/api/admin/login';
        requestBody = {
            email: identifier,
            password: password
        };
    } else {
        showMessage(loginMessage, 'Invalid login type.', true);
        return;
    }

    try {
        showMessage(loginMessage, `Logging in as ${loginType}...`);
        const data = await apiRequest(endpoint, 'POST', requestBody);

        if (data && data.token) {
            // Login successful
            authToken = data.token;
            userRole = loginType; // Set role based on how they logged in
            // Extract user ID from response (adjust keys based on your actual backend response)
            userId = (loginType === 'cricketer' && data.cricketer) ? data.cricketer.id :
                   (loginType === 'admin' && data.admin) ? data.admin.id : null;

            localStorage.setItem('authToken', authToken);
            localStorage.setItem('userRole', userRole);
            if (userId) {
                 localStorage.setItem('userId', userId);
            }

            showMessage(loginMessage, 'Login successful!');
            identifierInput.value = ''; // Clear form
            passwordInput.value = '';
            checkLoginState(); // Update UI
        } else {
            throw new Error('Login failed: No token received.');
        }
    } catch (error) {
        showMessage(loginMessage, `Login failed: ${error.message}`, true);
        // Clear potentially bad stored token if login fails
        localStorage.removeItem('authToken');
        localStorage.removeItem('userRole');
        localStorage.removeItem('userId');
        authToken = null;
        userRole = null;
        userId = null;
    }
}

async function handleSignup(event) {
    event.preventDefault();
    showMessage(signupMessage, ''); // Clear previous messages

    const nameInput = document.getElementById('signup-name');
    const emailInput = document.getElementById('signup-email');
    const mobileInput = document.getElementById('signup-mobile');
    const passwordInput = document.getElementById('signup-password');

    const name = nameInput.value.trim();
    const email = emailInput.value.trim();
    const mobile = mobileInput.value.trim();
    const password = passwordInput.value.trim();

    if (!name || !email || !mobile || !password) {
        showMessage(signupMessage, 'All fields are required.', true);
        return;
    }

    const requestBody = {
        name: name,
        email: email,
        mobile: mobile,
        password: password
    };

    try {
        showMessage(signupMessage, 'Signing up...');
        // The signup endpoint doesn't require authentication
        const data = await apiRequest('/api/signup', 'POST', requestBody, false);

        showMessage(signupMessage, data?.message || 'Signup successful! Please log in.');
        // Clear the form
        nameInput.value = '';
        emailInput.value = '';
        mobileInput.value = '';
        passwordInput.value = '';

    } catch (error) {
        showMessage(signupMessage, `Signup failed: ${error.message}`, true);
    }
}

function handleLogout() {
    // No API call needed for logout with JWT on the client-side
    // We just delete the token locally.
    showMessage(loginMessage, 'Logged out successfully.'); // Show message in login area

    // Clear state variables
    authToken = null;
    userRole = null;
    userId = null;

    // Clear localStorage
    localStorage.removeItem('authToken');
    localStorage.removeItem('userRole');
    localStorage.removeItem('userId');

    // Update UI
    checkLoginState();

    // Optional: Clear dynamic content areas
    announcementsList.innerHTML = '';
    profileDetails.innerHTML = '';
    cricketersList.innerHTML = '';
    showMessage(announcementsMessage, '');
    showMessage(profileMessage, '');
    showMessage(cricketersMessage, '');
    showMessage(createAnnouncementMessage, '');
}

async function handleUpdateProfile(event) {
    event.preventDefault();
    if (!authToken || userRole !== 'cricketer') return;

    showMessage(profileMessage, ''); // Clear previous messages

    const nameInput = document.getElementById('update-name');
    const mobileInput = document.getElementById('update-mobile');
    const passwordInput = document.getElementById('update-password');

    const name = nameInput.value.trim();
    const mobile = mobileInput.value.trim();
    const password = passwordInput.value.trim(); // Get password value

    const requestBody = {};
    if (name) requestBody.name = name;
    if (mobile) requestBody.mobile = mobile;
    if (password) requestBody.password = password; // Include password only if provided

    if (Object.keys(requestBody).length === 0) {
        showMessage(profileMessage, 'No changes detected to update.', true);
        return;
    }

    try {
        showMessage(profileMessage, 'Updating profile...');
        const data = await apiRequest('/api/cricketer/profile', 'PUT', requestBody, true);

        showMessage(profileMessage, 'Profile updated successfully!');
        passwordInput.value = ''; // Clear password field after update
        // Optionally re-fetch profile details to show updated data
        fetchCricketerProfile();

    } catch (error) {
        showMessage(profileMessage, `Profile update failed: ${error.message}`, true);
    }
}

async function handleCreateAnnouncement(event) {
    event.preventDefault();
    if (!authToken || userRole !== 'admin') return;

    showMessage(createAnnouncementMessage, ''); // Clear previous messages

    const titleInput = document.getElementById('announcement-title');
    const contentInput = document.getElementById('announcement-content');

    const title = titleInput.value.trim();
    const content = contentInput.value.trim();

    if (!title || !content) {
        showMessage(createAnnouncementMessage, 'Title and content are required.', true);
        return;
    }

    const requestBody = {
        title: title,
        content: content
    };

    try {
        showMessage(createAnnouncementMessage, 'Creating announcement...');
        const data = await apiRequest('/api/admin/announcements', 'POST', requestBody, true);

        showMessage(createAnnouncementMessage, 'Announcement created successfully!');
        titleInput.value = ''; // Clear form
        contentInput.value = '';
        // Re-fetch announcements to show the new one
        fetchAnnouncements();

    } catch (error) {
        showMessage(createAnnouncementMessage, `Failed to create announcement: ${error.message}`, true);
    }
}

async function fetchAnnouncements() {
    if (!authToken) return; // Need to be logged in

    showMessage(announcementsMessage, 'Loading announcements...');
    try {
        const data = await apiRequest('/api/announcements', 'GET', null, true);
        announcementsList.innerHTML = ''; // Clear previous list

        if (data && data.length > 0) {
            data.forEach(announcement => {
                const li = document.createElement('li');
                // Format date for better readability
                const date = new Date(announcement.createdAt).toLocaleString();
                li.innerHTML = `<strong>${announcement.title || 'Untitled'}</strong> (by ${announcement.createdBy}, ${date})<br>${announcement.content}`;
                announcementsList.appendChild(li);
            });
            showMessage(announcementsMessage, ''); // Clear loading message
        } else {
            showMessage(announcementsMessage, 'No announcements found.');
        }
    } catch (error) {
        showMessage(announcementsMessage, `Error fetching announcements: ${error.message}`, true);
    }
}

async function fetchCricketerProfile() {
    if (!authToken || userRole !== 'cricketer') return;

    showMessage(profileMessage, 'Loading profile...');
    try {
        const data = await apiRequest('/api/cricketer/profile', 'GET', null, true);
        profileDetails.innerHTML = ''; // Clear previous details

        if (data) {
            profileDetails.innerHTML = `
                <p><strong>ID:</strong> ${data.id}</p>
                <p><strong>Name:</strong> ${data.name}</p>
                <p><strong>Email:</strong> ${data.email}</p>
                <p><strong>Mobile:</strong> ${data.mobile}</p>
            `;
            // Pre-fill update form with current data
            document.getElementById('update-name').value = data.name || '';
            document.getElementById('update-mobile').value = data.mobile || '';
            showMessage(profileMessage, ''); // Clear loading message
        } else {
            showMessage(profileMessage, 'Could not load profile data.', true);
        }
    } catch (error) {
        showMessage(profileMessage, `Error fetching profile: ${error.message}`, true);
    }
}

async function fetchAllCricketers() {
    if (!authToken || userRole !== 'admin') return;

    showMessage(cricketersMessage, 'Loading cricketers list...');
    try {
        const data = await apiRequest('/api/admin/cricketers', 'GET', null, true);
        cricketersList.innerHTML = ''; // Clear previous list

        if (data && data.length > 0) {
            data.forEach(cricketer => {
                const li = document.createElement('li');
                li.innerHTML = `<strong>${cricketer.name}</strong> (ID: ${cricketer.id})<br>Email: ${cricketer.email}, Mobile: ${cricketer.mobile}`;
                cricketersList.appendChild(li);
            });
            showMessage(cricketersMessage, ''); // Clear loading message
        } else {
            showMessage(cricketersMessage, 'No cricketers found.');
        }
    } catch (error) {
        showMessage(cricketersMessage, `Error fetching cricketers: ${error.message}`, true);
    }
}

// --- Add more functions as needed --- 
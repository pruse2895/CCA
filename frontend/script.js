// DOM Elements
const loginPage = document.getElementById('login-page');
const signupPage = document.getElementById('signup-page');
const adminLoginPage = document.getElementById('admin-login-page');
const loggedInSection = document.getElementById('logged-in-section');
const showSignupLink = document.getElementById('show-signup');
const showLoginLink = document.getElementById('show-login');
const showAdminLoginLink = document.getElementById('show-admin-login');
const logoutBtn = document.getElementById('logout-btn');

// Form elements
const loginForm = document.getElementById('login-form');
const signupForm = document.getElementById('signup-form');
const announcementForm = document.getElementById('announcement-form');
const adminLoginForm = document.getElementById('admin-login-form');

// Message elements
const loginMessage = document.getElementById('login-message');
const signupMessage = document.getElementById('signup-message');

// List elements
const announcementsList = document.getElementById('announcements-list');
const cricketersList = document.getElementById('cricketers-list');
const profileInfo = document.getElementById('profile-info');

// Page navigation
function showPage(page) {
    if (!page) {
        console.error('Page element not found');
        return;
    }
    
    // Hide all pages
    document.querySelectorAll('.page').forEach(p => {
        p.style.display = 'none';
        p.classList.remove('active');
    });
    
    // Show the requested page
    page.style.display = 'flex';
    page.classList.add('active');
}

// Initialize event listeners
function initializeEventListeners() {
    console.log('Initializing event listeners...');
    
    // Signup link event listener
    if (showSignupLink) {
        showSignupLink.addEventListener('click', (e) => {
            e.preventDefault();
            console.log('Signup link clicked');
            showPage(signupPage);
        });
    } else {
        console.error('Signup link not found');
    }

    // Login link event listener
    if (showLoginLink) {
        showLoginLink.addEventListener('click', (e) => {
            e.preventDefault();
            console.log('Login link clicked');
            showPage(loginPage);
        });
    } else {
        console.error('Login link not found');
    }

    // Admin login link event listener
    if (showAdminLoginLink) {
        showAdminLoginLink.addEventListener('click', (e) => {
            e.preventDefault();
            showPage(adminLoginPage);
        });
    }

    // Signup form submission
    if (signupForm) {
        signupForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            console.log('Signup form submitted');
            
            // Get form values
            const name = document.getElementById('signup-name').value;
            const email = document.getElementById('signup-email').value;
            const mobile = document.getElementById('signup-mobile').value;
            const password = document.getElementById('signup-password').value;

            console.log('Attempting signup with:', { name, email, mobile });

            try {
                // First check if server is reachable
                try {
                    console.log('Checking server health...');
                    const healthCheck = await fetch('http://localhost:8080/', {
                        method: 'GET',
                        headers: {
                            'Accept': 'application/json',
                        },
                    });
                    console.log('Health check response status:', healthCheck.status);
                    
                    if (!healthCheck.ok) {
                        throw new Error(`Server returned status ${healthCheck.status}`);
                    }
                    
                    const healthData = await healthCheck.text();
                    console.log('Health check response:', healthData);
                } catch (healthError) {
                    console.error('Server health check failed:', healthError);
                    if (signupMessage) {
                        signupMessage.textContent = `Server connection failed: ${healthError.message}. Please ensure the server is running on port 8080.`;
                        signupMessage.classList.remove('success');
                        signupMessage.classList.add('error');
                    }
                    return;
                }

                // If health check passed, proceed with signup
                console.log('Server is healthy, proceeding with signup...');
                const response = await fetch('http://localhost:8080/api/signup', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name,
                        email,
                        mobile,
                        password
                    }),
                });

                console.log('Signup response status:', response.status);
                
                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Signup failed:', errorData);
                    if (signupMessage) {
                        signupMessage.textContent = errorData.error || 'Signup failed. Please try again.';
                        signupMessage.classList.remove('success');
                        signupMessage.classList.add('error');
                    }
                    return;
                }

                const data = await response.json();
                console.log('Signup response data:', data);

                if (data.message) {
                    signupMessage.textContent = 'Signup successful! Please login.';
                    signupMessage.classList.remove('error');
                    signupMessage.classList.add('success');
                    
                    // Clear form
                    signupForm.reset();
                    
                    // Show login page after successful signup
                    setTimeout(() => {
                        showPage(loginPage);
                    }, 1500);
                } else {
                    console.error('Invalid response format:', data);
                    if (signupMessage) {
                        signupMessage.textContent = 'Invalid response from server';
                        signupMessage.classList.remove('success');
                        signupMessage.classList.add('error');
                    }
                }
            } catch (error) {
                console.error('Signup error:', error);
                if (signupMessage) {
                    if (error.message.includes('Failed to fetch')) {
                        signupMessage.textContent = 'Cannot connect to server. Please check if the server is running and accessible.';
                    } else {
                        signupMessage.textContent = `Error: ${error.message}`;
                    }
                    signupMessage.classList.remove('success');
                    signupMessage.classList.add('error');
                }
            }
        });
    } else {
        console.error('Signup form not found');
    }

    // Login form submission
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            console.log('Login form submitted');
            
            const mobile = document.getElementById('login-mobile').value;
            const password = document.getElementById('login-password').value;

            console.log('Attempting login with mobile:', mobile);

            try {
                const response = await fetch('http://localhost:8080/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ mobile, password }),
                });

                console.log('Login response status:', response.status);
                
                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Login failed:', errorData);
                    if (loginMessage) {
                        loginMessage.textContent = errorData.error || 'Invalid mobile number or password';
                        loginMessage.classList.add('error');
                    }
                    return;
                }

                const data = await response.json();
                console.log('Login response data:', data);

                if (data.token && data.cricketer) {
                    console.log('Login successful, storing token and user data');
                    localStorage.setItem('token', data.token);
                    localStorage.setItem('user', JSON.stringify(data.cricketer));
                    
                    // Clear any previous error messages
                    if (loginMessage) {
                        loginMessage.textContent = '';
                        loginMessage.classList.remove('error');
                    }
                    
                    // Show logged in section
                    showPage(loggedInSection);
                    loadUserData();
                } else {
                    console.error('Invalid response format:', data);
                    if (loginMessage) {
                        loginMessage.textContent = 'Invalid response from server';
                        loginMessage.classList.add('error');
                    }
                }
            } catch (error) {
                console.error('Login error:', error);
                if (loginMessage) {
                    loginMessage.textContent = 'Failed to connect to server. Please try again.';
                    loginMessage.classList.add('error');
                }
            }
        });
    } else {
        console.error('Login form not found');
    }

    // Logout button
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            showPage(loginPage);
        });
    } else {
        console.error('Logout button not found');
    }

    // Admin login form submission
    if (adminLoginForm) {
        adminLoginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            console.log('Admin login form submitted');
            
            const username = document.getElementById('admin-username').value;
            const password = document.getElementById('admin-password').value;

            console.log('Attempting admin login with username:', username);

            try {
                // First check if server is reachable
                try {
                    console.log('Checking server health...');
                    const healthCheck = await fetch('http://localhost:8080/', {
                        method: 'GET',
                        headers: {
                            'Accept': 'application/json',
                        },
                    });
                    console.log('Health check response status:', healthCheck.status);
                    
                    if (!healthCheck.ok) {
                        throw new Error(`Server returned status ${healthCheck.status}`);
                    }
                    
                    const healthData = await healthCheck.text();
                    console.log('Health check response:', healthData);
                } catch (healthError) {
                    console.error('Server health check failed:', healthError);
                    alert('Server connection failed. Please ensure the server is running on port 8080.');
                    return;
                }

                // If health check passed, proceed with admin login
                console.log('Server is healthy, proceeding with admin login...');
                const response = await fetch('http://localhost:8080/api/admin/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ username, password }),
                });

                console.log('Admin login response status:', response.status);
                
                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Admin login failed:', errorData);
                    alert(errorData.error || 'Invalid admin credentials');
                    return;
                }

                const data = await response.json();
                console.log('Admin login response data:', data);

                if (data.token) {
                    console.log('Admin login successful, storing token');
                    localStorage.setItem('token', data.token);
                    localStorage.setItem('user', JSON.stringify({ role: 'admin' }));
                    
                    // Show admin section
                    showPage(loggedInSection);
                    loadUserData();
                } else {
                    console.error('Invalid response format:', data);
                    alert('Invalid response from server');
                }
            } catch (error) {
                console.error('Admin login error:', error);
                if (error.message.includes('Failed to fetch')) {
                    alert('Cannot connect to server. Please check if the server is running.');
                } else {
                    alert(`Error: ${error.message}`);
                }
            }
        });
    }

    // Announcement form
    if (announcementForm) {
        announcementForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            console.log('Announcement form submitted');
            
            const token = localStorage.getItem('token');
            if (!token) {
                console.error('No token found');
                alert('Please login first');
                return;
            }

            // Verify admin role
            const user = JSON.parse(localStorage.getItem('user'));
            if (!user || user.role !== 'admin') {
                console.error('User is not an admin');
                alert('Only admins can create announcements');
                return;
            }

            const title = document.getElementById('announcement-title').value;
            const content = document.getElementById('announcement-content').value;

            console.log('Creating announcement with:', { title, content });

            try {
                const response = await fetch('http://localhost:8080/api/admin/announcements', {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ title, content }),
                });

                console.log('Announcement response status:', response.status);
                
                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Announcement creation failed:', errorData);
                    alert(errorData.error || 'Failed to create announcement. Please check if you have admin privileges.');
                    return;
                }

                const data = await response.json();
                console.log('Announcement created successfully:', data);

                // Clear form
                announcementForm.reset();
                
                // Reload announcements
                loadUserData();
                
                // Show success message
                alert('Announcement created successfully!');
            } catch (error) {
                console.error('Error creating announcement:', error);
                if (error.message.includes('Failed to fetch')) {
                    alert('Cannot connect to server. Please check if the server is running.');
                } else {
                    alert(`Error: ${error.message}`);
                }
            }
        });
    }
}

// Load user data
async function loadUserData() {
    const token = localStorage.getItem('token');
    if (!token) {
        showPage(loginPage);
        return;
    }

    try {
        // Load profile
        const profileResponse = await fetch('/api/cricketer/profile', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (profileResponse.ok) {
            const profile = await profileResponse.json();
            displayProfile(profile);
        }

        // Load announcements
        const announcementsResponse = await fetch('/api/announcements', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (announcementsResponse.ok) {
            const announcements = await announcementsResponse.json();
            displayAnnouncements(announcements);
        }

        // Check if user is admin
        const user = JSON.parse(localStorage.getItem('user'));
        if (user && user.role === 'admin') {
            document.getElementById('admin-section').style.display = 'block';
            loadCricketers();
        } else {
            document.getElementById('admin-section').style.display = 'none';
        }
    } catch (error) {
        console.error('Error loading user data:', error);
    }
}

// Display profile
function displayProfile(profile) {
    if (!profileInfo) return;
    profileInfo.innerHTML = `
        <p><strong>Name:</strong> ${profile.name}</p>
        <p><strong>Email:</strong> ${profile.email}</p>
        <p><strong>Mobile:</strong> ${profile.mobile}</p>
    `;
}

// Display announcements
function displayAnnouncements(announcements) {
    if (!announcementsList) return;
    announcementsList.innerHTML = announcements.map(announcement => `
        <li>
            <h3>${announcement.title}</h3>
            <p>${announcement.content}</p>
            <small>${new Date(announcement.createdAt).toLocaleString()}</small>
        </li>
    `).join('');
}

// Load cricketers (admin only)
async function loadCricketers() {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch('/api/admin/cricketers', {
            headers: {
                'Authorization': `Bearer ${token}`,
            },
        });

        if (response.ok) {
            const cricketers = await response.json();
            if (cricketersList) {
                cricketersList.innerHTML = cricketers.map(cricketer => `
                    <li>
                        <p><strong>Name:</strong> ${cricketer.name}</p>
                        <p><strong>Email:</strong> ${cricketer.email}</p>
                        <p><strong>Mobile:</strong> ${cricketer.mobile}</p>
                    </li>
                `).join('');
            }
        }
    } catch (error) {
        console.error('Error loading cricketers:', error);
    }
}

// Initialize everything when the DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, initializing...');
    
    // Initialize event listeners
    initializeEventListeners();
    
    // Check authentication on page load
    const token = localStorage.getItem('token');
    if (token) {
        showPage(loggedInSection);
        loadUserData();
    } else {
        showPage(loginPage);
    }
});

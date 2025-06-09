// API functions
const api = {
    async getUsers() {
        const response = await fetch('/api/users');
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        const data = await response.json();
        return data.users;
    },
    async createUser(user) {
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(user)
        });
        if (!response.ok) {
            throw new Error('Failed to create user');
        }
        return response.json();
    }
};

// Display users in the UI
function displayUsers(users) {
    const usersList = document.getElementById('usersList');
    usersList.innerHTML = '';
    
    if (users.length === 0) {
        usersList.innerHTML = '<p>No users found.</p>';
        return;
    }

    users.forEach(user => {
        const userCard = document.createElement('div');
        userCard.className = 'user-card';
        userCard.innerHTML = `
            <h3>${user.name}</h3>
            <p>ID: ${user.id}</p>
        `;
        usersList.appendChild(userCard);
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const createUserForm = document.getElementById('createUserForm');
    const getUsersBtn = document.getElementById('getUsersBtn');

    // Handle form submission
    createUserForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const user = {
            id: document.getElementById('userId').value,
            name: document.getElementById('userName').value
        };

        try {
            await api.createUser(user);
            alert('User created successfully!');
            createUserForm.reset();
            const users = await api.getUsers();
            displayUsers(users);
        } catch (error) {
            alert('Error creating user: ' + error.message);
        }
    });

    // Handle get users button click
    getUsersBtn.addEventListener('click', async () => {
        try {
            const users = await api.getUsers();
            displayUsers(users);
        } catch (error) {
            alert('Error fetching users: ' + error.message);
        }
    });
});

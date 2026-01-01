// Device detail JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on a device detail page by examining the URL
    const path = window.location.pathname;
    
    // If we're at a device path, load the specific device
    if (path.startsWith('/device/')) {
        const deviceId = path.split('/').pop();
        if (deviceId) {
            loadDeviceDetails(deviceId);
        } else {
            document.getElementById('deviceDetails').innerHTML = 
                '<p>Invalid device ID</p>';
        }
    } else {
        // If we're not on a device page, show a message
        document.getElementById('deviceDetails').innerHTML = 
            '<p>This is the device details page. Please navigate here from the dashboard.</p>';
    }

    // Function to load device details
    async function loadDeviceDetails(id) {
        try {
            const response = await fetch(`/api/device/${id}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const device = await response.json();
            displayDeviceDetails(device);
        } catch (error) {
            console.error('Error loading device details:', error);
            document.getElementById('deviceDetails').innerHTML = 
                '<p>Error loading device details</p>';
        }
    }

    // Function to display device details
    function displayDeviceDetails(device) {
        const container = document.getElementById('deviceDetails');
        
        let txtRecordsHtml = '';
        if (device.txt && Object.keys(device.txt).length > 0) {
            txtRecordsHtml = '<div class="txt-container">' + 
                Object.entries(device.txt).map(([key, value]) => 
                    `<div class="txt-record">${escapeHtml(key)}=${escapeHtml(value)}</div>`
                ).join('') + '</div>';
        } else {
            txtRecordsHtml = '<p>No TXT records available</p>';
        }

        const addressesHtml = device.addresses && device.addresses.length > 0 
            ? device.addresses.map(addr => `<div class="address-item">${escapeHtml(addr)}</div>`).join('')
            : '<p>No addresses available</p>';

        container.innerHTML = `
            <h2>Device Details</h2>
            <div class="detail-row">
                <div class="detail-label">ID:</div>
                <div class="detail-value">${escapeHtml(device.id)}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Name:</div>
                <div class="detail-value">${escapeHtml(device.name)}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Host:</div>
                <div class="detail-value">${escapeHtml(device.host)}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Port:</div>
                <div class="detail-value">${device.port}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Protocol:</div>
                <div class="detail-value">${escapeHtml(device.protocol || 'N/A')}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Service:</div>
                <div class="detail-value">${escapeHtml(device.service)}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Addresses:</div>
                <div class="detail-value">${addressesHtml}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">TTL:</div>
                <div class="detail-value">${device.ttl}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">Last Seen:</div>
                <div class="detail-value">${formatDate(device.last_seen)}</div>
            </div>
            <div class="detail-row">
                <div class="detail-label">TXT Records:</div>
                <div class="detail-value">${txtRecordsHtml}</div>
            </div>
        `;
    }

    // Helper function to format date
    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleString();
    }

    // Helper function to escape HTML
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
});
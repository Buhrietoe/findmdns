// Main dashboard JavaScript
document.addEventListener('DOMContentLoaded', function() {
    const scanBtn = document.getElementById('scanBtn');
    const devicesBody = document.getElementById('devicesBody');
    const lastScanSpan = document.getElementById('lastScan');

    // Check current URL to determine if we should show dashboard or device detail
    const path = window.location.pathname;
    
    if (path.startsWith('/device/')) {
        // Show device detail view
        const deviceId = path.split('/').pop();
        if (deviceId) {
            // Decode the device ID if it was URL encoded
            const decodedId = decodeURIComponent(deviceId);
            loadDeviceDetails(decodedId);
        } else {
            document.getElementById('devicesTableContainer').innerHTML = '<p>Invalid device ID</p>';
        }
    } else {
        // Show dashboard view
        loadDevices();
    }

    // Scan button event listener
    scanBtn.addEventListener('click', function() {
        scanDevices();
    });

    // Function to load devices from API
    async function loadDevices() {
        try {
            const response = await fetch('/api/devices');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const devices = await response.json();
            displayDevices(devices);
            updateLastScan();
        } catch (error) {
            console.error('Error loading devices:', error);
            devicesBody.innerHTML = '<tr><td colspan="6">Error loading devices</td></tr>';
        }
    }

    // Function to scan for new devices
    async function scanDevices() {
        try {
            const response = await fetch('/api/scan', { method: 'POST' });
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            console.log('Scan started:', result);
            
            // Refresh the device list after a short delay
            setTimeout(loadDevices, 1000);
        } catch (error) {
            console.error('Error starting scan:', error);
        }
    }

    // Function to update last scan time
    async function updateLastScan() {
        try {
            const response = await fetch('/api/last-scan');
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            const lastScanDate = new Date(result.last_scan);
            lastScanSpan.textContent = `Last scan: ${lastScanDate.toLocaleString()}`;
        } catch (error) {
            console.error('Error getting last scan time:', error);
            lastScanSpan.textContent = 'Last scan: Unknown';
        }
    }

    // Function to display devices in the table
    function displayDevices(devices) {
        if (devices.length === 0) {
            devicesBody.innerHTML = '<tr><td colspan="6">No devices found</td></tr>';
            return;
        }

        const rows = devices.map(device => `
            <tr>
                <td>${escapeHtml(device.name)}</td>
                <td>${escapeHtml(device.host)}</td>
                <td>${device.port}</td>
                <td>${formatAddresses(device.addresses)}</td>
                <td>${formatDate(device.last_seen)}</td>
                <td><a href="/device/${encodeURIComponent(device.id)}" class="device-link">View Details</a></td>
            </tr>
        `).join('');

        devicesBody.innerHTML = rows;
    }

    // Function to load and display device details
    async function loadDeviceDetails(id) {
        try {
            const response = await fetch(`/api/device/${encodeURIComponent(id)}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const device = await response.json();
            displayDeviceDetails(device);
        } catch (error) {
            console.error('Error loading device details:', error);
            document.getElementById('devicesTableContainer').innerHTML = 
                `<p>Error loading device details: ${error.message}</p>
                 <p>Device ID: ${escapeHtml(id)}</p>
                 <a href="/" class="device-link">← Back to Dashboard</a>`;
        }
    }

    // Function to display single device details
    function displayDeviceDetails(device) {
        const container = document.getElementById('devicesTableContainer');
        
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
            <h2>Device Details: ${escapeHtml(device.name)}</h2>
            <div class="detail-row">
                <div class="detail-label">ID:</div>
                <div class="detail-value">${escapeHtml(device.id)}</div>
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
            <div style="margin-top: 20px;">
                <a href="/" class="device-link">← Back to Dashboard</a>
            </div>
        `;
    }

    // Helper function to format addresses
    function formatAddresses(addresses) {
        if (!addresses || addresses.length === 0) {
            return 'N/A';
        }
        
        return `<div class="addresses-list">
            ${addresses.map(addr => `<span class="address-item">${escapeHtml(addr)}</span>`).join('')}
        </div>`;
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
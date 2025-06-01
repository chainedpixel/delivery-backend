
class DeliveryTracker {
    constructor(config) {
        this.config = config;
        this.map = null;
        this.driverMarker = null;
        this.pickupMarker = null;
        this.deliveryMarker = null;
        this.tracker = null;
        this.currentOrderId = null;
        this.orderData = {};
        
        mapboxgl.accessToken = this.config.MAPBOX_TOKEN;
    }
    
    async init() {
        this.currentOrderId = this.getOrderIdFromUrl();

        if (!this.currentOrderId) {
            this.showError('No se proporcionó un ID de pedido válido en la URL. Use: ?order_id=su-id-aqui');
            return;
        }

        await this.initializeApp();
    }
    
    getOrderIdFromUrl() {
        const urlParams = new URLSearchParams(window.location.search);
        return urlParams.get('order_id');
    }
    
    showError(message) {
        const errorContainer = document.getElementById('error-container');
        errorContainer.innerHTML = `
            <div class="error-message">
                <strong>Error:</strong> ${message}
                <br><br>
                <small>Ejemplo de URL válida: ${window.location.origin}${window.location.pathname}?order_id=eb32124a-6664-4083-a335-ea810ef7420e</small>
            </div>
        `;
    }
    
    showLoading(show = true) {
        document.getElementById('loading').classList.toggle('hidden', !show);
        document.getElementById('main-content').classList.toggle('hidden', show);
    }
    
    async apiRequest(endpoint, options = {}) {
        const url = `${this.config.API_BASE_URL}${endpoint}`;
        const defaultOptions = {
            headers: {
                'Authorization': `Bearer ${this.config.AUTH_TOKEN}`,
                'Content-Type': 'application/json'
            }
        };

        try {
            const response = await fetch(url, { ...defaultOptions, ...options });
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error('API Request failed:', error);
            throw error;
        }
    }
    
    async getOrderDetails(orderId) {
        return await this.apiRequest(`/orders/${orderId}`);
    }

    async startSimulation(orderId) {
        return await this.apiRequest(`/orders/${orderId}/simulate`, { method: 'POST' });
    }

    async assignRandomDriver(orderId) {
        return await this.apiRequest(`/orders/${orderId}/assign-driver`, { method: 'POST' });
    }

    async simulateMovement(orderId) {
        return await this.apiRequest(`/orders/${orderId}/simulate-movement`, { method: 'POST' });
    }
    
    async initializeApp() {
        this.showLoading(true);

        try {
            const response = await this.getOrderDetails(this.currentOrderId);
            this.orderData = response;
            
            this.updateOrderUI(this.orderData);
            
            this.initMap();
            
            this.initTracker();
            
            this.setupSimulationControls();

            this.showLoading(false);

        } catch (error) {
            this.showError(`Error al cargar el pedido: ${error.message}`);
            this.showLoading(false);
        }
    }
    
    updateOrderUI(order) {
        console.log('Updating UI with order data:', order);

        this.setElementText('order-id', order.id || '-');
        this.setElementText('tracking-number', order.tracking_number || '-');
        this.setElementText('client-name', order.client?.full_name || '-');
        this.setElementText('created-date', order.created_at ?
            new Date(order.created_at).toLocaleString('es-ES') : '-');
        this.setElementText('order-status', order.status || '-');
        
        this.updateAddresses(order);
        
        this.updateDriverInfo(order);
        
        this.updateStatusSteps(order.status);
        this.updateProgress(order.status);
        this.updateStatusDescription(order.status);
        
        document.title = `Rastreo - ${order.tracking_number || order.id}`;
    }
    
    updateAddresses(order) {
        console.log('Updating addresses with:', order);

        let pickupAddress = '-';
        let deliveryAddress = '-';
        
        if (order.detail) {
            pickupAddress = order.detail.pickup_address || pickupAddress;
            deliveryAddress = order.detail.delivery_address || deliveryAddress;
        }
        
        if (pickupAddress === '-' && order.pickup_address) {
            pickupAddress = this.formatAddress(order.pickup_address);
        }
        if (deliveryAddress === '-' && order.delivery_address) {
            deliveryAddress = this.formatAddress(order.delivery_address);
        }
        
        if (pickupAddress === '-' && order.PickupAddress) {
            pickupAddress = this.formatAddress(order.PickupAddress);
        }
        if (deliveryAddress === '-' && order.DeliveryAddress) {
            deliveryAddress = this.formatAddress(order.DeliveryAddress);
        }

        console.log('Final addresses:', { pickupAddress, deliveryAddress });

        this.setElementText('pickup-address', pickupAddress);
        this.setElementText('delivery-address', deliveryAddress);
    }
    
    formatAddress(addressObj) {
        if (!addressObj) return '-';

        if (typeof addressObj === 'string') {
            return addressObj;
        }

        const parts = [];
        if (addressObj.address_line1) parts.push(addressObj.address_line1);
        if (addressObj.address_line2) parts.push(addressObj.address_line2);
        if (addressObj.city) parts.push(addressObj.city);
        if (addressObj.state) parts.push(addressObj.state);

        return parts.length > 0 ? parts.join(', ') : '-';
    }
    
    updateDriverInfo(order) {
        console.log('Updating driver info with:', order);

        let driverName = 'Sin asignar';
        let vehicleInfo = '-';
        
        if (order.driver?.user) {
            driverName = order.driver.user.full_name || 'Conductor asignado';
            if (order.driver.vehicle_model || order.driver.vehicle_color) {
                vehicleInfo = `${order.driver.vehicle_model || ''} ${order.driver.vehicle_color || ''}`.trim();
            }
        }
        else if (order.Driver?.User) {
            driverName = order.Driver.User.full_name || order.Driver.User.FullName || 'Conductor asignado';
            if (order.Driver.vehicle_model || order.Driver.VehicleModel) {
                const model = order.Driver.vehicle_model || order.Driver.VehicleModel || '';
                const color = order.Driver.vehicle_color || order.Driver.VehicleColor || '';
                vehicleInfo = `${model} ${color}`.trim();
            }
        }
        else if (order.driver_id || order.DriverID) {
            driverName = 'Conductor asignado';
            vehicleInfo = 'Información no disponible';
        }

        console.log('Final driver info:', { driverName, vehicleInfo });

        this.setElementText('driver-name', driverName);
        this.setElementText('vehicle-info', vehicleInfo);
    }
    
    setElementText(id, text) {
        const element = document.getElementById(id);
        if (element) {
            element.textContent = text;
        }
    }
    
    updateStatusSteps(status) {
        const steps = this.config.UI.statusSteps;
        const currentIndex = steps.indexOf(status);

        document.querySelectorAll('.status-step').forEach((step, index) => {
            if (index <= currentIndex) {
                step.classList.add('active');
            } else {
                step.classList.remove('active');
            }
        });
    }
    
    updateProgress(status) {
        const progress = this.config.UI.progressMap[status] || 0;
        const progressElement = document.getElementById('order-progress');
        if (progressElement) {
            progressElement.style.width = `${progress}%`;
        }
    }
    
    updateStatusDescription(status) {
        const description = this.config.UI.statusDescriptions[status] || 'Estado desconocido';
        this.setElementText('status-description', description);
    }
    
    initMap() {
        const defaultLat = 13.6929;
        const defaultLng = -89.2182;

        this.map = new mapboxgl.Map({
            container: 'map',
            style: 'mapbox://styles/mapbox/streets-v12',
            center: [defaultLng, defaultLat],
            zoom: 13
        });

        this.map.on('load', () => {
            this.createMarkers();
            if (this.orderData.detail) {
                this.addRoute();
            }
        });
    }
    
    createMarkers() {
        const defaultCoords = {
            pickup: { lat: 13.6762, lng: -89.2874 },
            delivery: { lat: 13.6783, lng: -89.2353 },
            driver: { lat: 13.6772, lng: -89.2650 }
        };
        
        const driverEl = document.createElement('div');
        driverEl.className = 'marker marker-driver';
        this.driverMarker = new mapboxgl.Marker(driverEl)
            .setLngLat([defaultCoords.driver.lng, defaultCoords.driver.lat])
            .addTo(this.map);
        
        const pickupEl = document.createElement('div');
        pickupEl.className = 'marker marker-pickup';
        this.pickupMarker = new mapboxgl.Marker(pickupEl)
            .setLngLat([defaultCoords.pickup.lng, defaultCoords.pickup.lat])
            .setPopup(new mapboxgl.Popup().setText('Punto de Recogida'))
            .addTo(this.map);
        
        const deliveryEl = document.createElement('div');
        deliveryEl.className = 'marker marker-delivery';
        this.deliveryMarker = new mapboxgl.Marker(deliveryEl)
            .setLngLat([defaultCoords.delivery.lng, defaultCoords.delivery.lat])
            .setPopup(new mapboxgl.Popup().setText('Punto de Entrega'))
            .addTo(this.map);
    }
    
    async addRoute() {
        const pickupCoords = [this.pickupMarker.getLngLat().lng, this.pickupMarker.getLngLat().lat];
        const deliveryCoords = [this.deliveryMarker.getLngLat().lng, this.deliveryMarker.getLngLat().lat];

        try {
            const query = await fetch(
                `https://api.mapbox.com/directions/v5/mapbox/driving/${pickupCoords[0]},${pickupCoords[1]};${deliveryCoords[0]},${deliveryCoords[1]}?steps=true&geometries=geojson&access_token=${mapboxgl.accessToken}`
            );

            const json = await query.json();
            const route = json.routes[0].geometry.coordinates;

            if (this.map.getSource('route')) {
                this.map.getSource('route').setData({
                    type: 'Feature',
                    properties: {},
                    geometry: {
                        type: 'LineString',
                        coordinates: route
                    }
                });
            } else {
                this.map.addSource('route', {
                    type: 'geojson',
                    data: {
                        type: 'Feature',
                        properties: {},
                        geometry: {
                            type: 'LineString',
                            coordinates: route
                        }
                    }
                });

                this.map.addLayer({
                    id: 'route',
                    type: 'line',
                    source: 'route',
                    layout: {
                        'line-join': 'round',
                        'line-cap': 'round'
                    },
                    paint: {
                        'line-color': '#4CAF50',
                        'line-width': 5,
                        'line-opacity': 0.8
                    }
                });
            }
            
            const bounds = new mapboxgl.LngLatBounds()
                .extend(pickupCoords)
                .extend(deliveryCoords);
            this.map.fitBounds(bounds, { padding: 60 });

        } catch (error) {
            console.error('Error al cargar la ruta:', error);
        }
    }
    
    initTracker() {
        this.tracker = new OrderTracker(this.config.AUTH_TOKEN, this.currentOrderId, this.config.WS_BASE_URL);

        this.tracker.onOpen = () => {
            console.log('Conexión WebSocket establecida');
            this.updateConnectionStatus('Conectado');
        };

        this.tracker.onClose = () => {
            console.log('Conexión WebSocket cerrada');
            this.updateConnectionStatus('Desconectado');
        };

        this.tracker.onError = (error) => {
            console.error('Error en WebSocket:', error);
            this.updateConnectionStatus(`Error: ${error.message}`);
        };

        this.tracker.onOrderUpdate = (data) => {
            console.log('Actualización de pedido recibida:', data);
            this.handleOrderUpdate(data);
        };

        this.tracker.onLocationUpdate = (data) => {
            console.log('Actualización de ubicación recibida:', data);
            this.handleLocationUpdate(data);
        };

        this.tracker.connect();
    }
    
    updateConnectionStatus(status) {
        const statusElement = document.getElementById('connection-status');
        if (!statusElement) return;

        statusElement.textContent = status;
        statusElement.className = 'status-badge';

        if (status === 'Conectado') {
            statusElement.classList.add('status-connected');
        } else if (status === 'Desconectado') {
            statusElement.classList.add('status-disconnected');
        } else if (status.includes('Error') || status.includes('Reconectando')) {
            statusElement.classList.add('status-reconnecting');
        }
    }
    
    handleOrderUpdate(data) {
        console.log('Processing order update:', data);
        
        this.setElementText('order-status', data.status);
        this.setElementText('status-description', data.description);
        
        this.updateStatusSteps(data.status);
        this.updateProgress(data.status);
        
        if (data.order) {
            if (data.order.company_name) {
                console.log('Company name:', data.order.company_name);
            }

            if (data.order.driver_name) {
                this.setElementText('driver-name', data.order.driver_name);
            } else if (data.order.driver_id) {
                this.setElementText('driver-name', 'Conductor asignado');
            }

            if (data.order.client_name) {
                this.setElementText('client-name', data.order.client_name);
            }

            if (data.order.tracking_number) {
                this.setElementText('tracking-number', data.order.tracking_number);
            }
            
            this.orderData = { ...this.orderData, ...data.order, status: data.status };
        }
        
        this.showSimulationMessage(data.description, 'success');
    }
    
    handleLocationUpdate(data) {
        this.setElementText('driver-location',
            `Lat: ${data.latitude.toFixed(6)}, Lng: ${data.longitude.toFixed(6)}`);
        
        if (this.map && this.driverMarker) {
            this.driverMarker.setLngLat([data.longitude, data.latitude]);
            this.map.panTo([data.longitude, data.latitude]);
        }
    }
    
    setupSimulationControls() {
        const startSimBtn = document.getElementById('start-simulation');
        if (startSimBtn) {
            startSimBtn.addEventListener('click', async () => {
                try {
                    this.disableAllSimulationButtons();
                    this.showSimulationMessage('Iniciando simulación completa...', 'success');

                    await this.startSimulation(this.currentOrderId);
                    this.showSimulationMessage('Simulación iniciada correctamente. El pedido progresará automáticamente.', 'success');
                    
                } catch (error) {
                    console.error('Error starting simulation:', error);
                    this.showSimulationMessage(`Error al iniciar simulación: ${this.getErrorMessage(error)}`, 'error');
                    this.enableAllSimulationButtons();
                }
            });
        }

        const assignDriverBtn = document.getElementById('assign-driver');
        if (assignDriverBtn) {
            assignDriverBtn.addEventListener('click', async () => {
                try {
                    this.setButtonLoading('assign-driver', 'Asignando...');
                    await this.assignRandomDriver(this.currentOrderId);
                    this.showSimulationMessage('Conductor asignado correctamente');
                    this.setButtonNormal('assign-driver', 'Asignar Conductor');
                } catch (error) {
                    console.error('Error assigning driver:', error);
                    if (error.message.includes('foreign key constraint fails') ||
                        error.message.includes('driver')) {
                        this.showSimulationMessage('No hay conductores disponibles. La simulación continuará sin conductor real.', 'warning');
                    } else {
                        this.showSimulationMessage(`Error al asignar conductor: ${this.getErrorMessage(error)}`, 'error');
                    }
                    this.setButtonNormal('assign-driver', 'Asignar Conductor');
                }
            });
        }

        const simulateMovBtn = document.getElementById('simulate-movement');
        if (simulateMovBtn) {
            simulateMovBtn.addEventListener('click', async () => {
                try {
                    this.setButtonLoading('simulate-movement', 'Simulando...');
                    await this.simulateMovement(this.currentOrderId);
                    this.showSimulationMessage('Simulación de movimiento iniciada');
                    this.setButtonNormal('simulate-movement', 'Simular Movimiento');
                } catch (error) {
                    console.error('Error simulating movement:', error);
                    this.showSimulationMessage(`Error al simular movimiento: ${this.getErrorMessage(error)}`, 'error');
                    this.setButtonNormal('simulate-movement', 'Simular Movimiento');
                }
            });
        }
    }
    
    disableAllSimulationButtons() {
        const buttons = ['start-simulation', 'assign-driver', 'simulate-movement'];
        buttons.forEach(buttonId => {
            const button = document.getElementById(buttonId);
            if (button) {
                button.disabled = true;
            }
        });
    }
    
    enableAllSimulationButtons() {
        const buttons = ['start-simulation', 'assign-driver', 'simulate-movement'];
        buttons.forEach(buttonId => {
            const button = document.getElementById(buttonId);
            if (button) {
                button.disabled = false;
            }
        });
    }
    
    setButtonLoading(buttonId, text) {
        const button = document.getElementById(buttonId);
        if (button) {
            button.disabled = true;
            button.dataset.originalText = button.textContent;
            button.textContent = text;
        }
    }
    
    setButtonNormal(buttonId, text) {
        const button = document.getElementById(buttonId);
        if (button) {
            button.disabled = false;
            button.textContent = text || button.dataset.originalText || button.textContent;
        }
    }
    
    getErrorMessage(error) {
        const errorStr = error.message || error.toString();

        if (errorStr.includes('foreign key constraint fails')) {
            return this.config.UI.errorMessages['DRIVER_ASSIGNMENT_ERROR'];
        } else if (errorStr.includes('not found')) {
            return this.config.UI.errorMessages['ORDER_NOT_FOUND'];
        } else if (errorStr.includes('network') || errorStr.includes('fetch')) {
            return this.config.UI.errorMessages['NETWORK_ERROR'];
        } else {
            return errorStr;
        }
    }
    
    showSimulationMessage(message, type = 'success') {
        const msgElement = document.createElement('div');

        let backgroundColor;
        switch(type) {
            case 'error':
                backgroundColor = '#dc3545';
                break;
            case 'warning':
                backgroundColor = '#ffc107';
                msgElement.style.color = '#212529';
                break;
            default:
                backgroundColor = '#28a745';
        }

        msgElement.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 10px 20px;
            border-radius: 4px;
            color: ${type === 'warning' ? '#212529' : 'white'};
            font-weight: 500;
            z-index: 10000;
            background-color: ${backgroundColor};
            box-shadow: 0 2px 8px rgba(0,0,0,0.15);
            max-width: 400px;
            word-wrap: break-word;
        `;
        msgElement.textContent = message;

        document.body.appendChild(msgElement);

        setTimeout(() => {
            if (document.body.contains(msgElement)) {
                msgElement.style.opacity = '0';
                msgElement.style.transition = 'opacity 0.3s ease';
                setTimeout(() => {
                    if (document.body.contains(msgElement)) {
                        document.body.removeChild(msgElement);
                    }
                }, 300);
            }
        }, 5000);
    }
    
    disableButton(buttonId) {
        const button = document.getElementById(buttonId);
        if (button) {
            button.disabled = true;
            button.textContent = 'En progreso...';
        }
    }
}

class OrderTracker {
    constructor(token, orderID, wsBaseUrl) {
        this.token = token;
        this.orderID = orderID;
        this.wsBaseUrl = wsBaseUrl;
        this.socket = null;
        this.connected = false;
        this.onOrderUpdate = null;
        this.onLocationUpdate = null;
        this.onError = null;
        this.onOpen = null;
        this.onClose = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectTimeout = null;
    }

    connect() {
        if (this.socket) {
            this.socket.close();
        }

        const socketUrl = `${this.wsBaseUrl}/tracking/ws?token=${encodeURIComponent(this.token)}`;
        console.log("Connecting to WebSocket:", socketUrl);

        this.socket = new WebSocket(socketUrl);

        this.socket.onopen = this.handleOpen.bind(this);
        this.socket.onmessage = this.handleMessage.bind(this);
        this.socket.onclose = this.handleClose.bind(this);
        this.socket.onerror = this.handleError.bind(this);
    }

    handleOpen(event) {
        console.log("WebSocket connection established!");
        this.connected = true;
        this.reconnectAttempts = 0;

        this.subscribeToOrder(this.orderID);

        if (this.onOpen) {
            this.onOpen(event);
        }
    }

    handleMessage(event) {
        try {
            console.log("WebSocket message received:", event.data);
            const message = JSON.parse(event.data);

            switch (message.type) {
                case 'ORDER_UPDATE':
                    if (this.onOrderUpdate) {
                        this.onOrderUpdate(message.data);
                    }
                    break;
                case 'LOCATION':
                    if (this.onLocationUpdate) {
                        this.onLocationUpdate(message.data);
                    }
                    break;
                case 'ERROR':
                    if (this.onError) {
                        this.onError(message.data);
                    }
                    break;
                default:
                    console.log("Unknown message type:", message.type);
            }
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    }

    handleClose(event) {
        this.connected = false;
        console.log("WebSocket connection closed:", event.code, event.reason);

        if (this.onClose) {
            this.onClose(event);
        }

        if (event.code !== 1000) {
            this.tryReconnect();
        }
    }

    handleError(error) {
        console.error("WebSocket connection error:", error);
        if (this.onError) {
            this.onError({
                code: 'CONNECTION_ERROR',
                message: 'Error en la conexión WebSocket',
                details: error.toString()
            });
        }
    }

    tryReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            if (this.onError) {
                this.onError({
                    code: 'RECONNECT_FAILED',
                    message: 'No se pudo reconectar después de varios intentos'
                });
            }
            return;
        }

        this.reconnectAttempts++;
        const delay = Math.min(30000, 1000 * Math.pow(2, this.reconnectAttempts));

        console.log(`Attempting to reconnect in ${delay/1000} seconds (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);

        this.reconnectTimeout = setTimeout(() => {
            if (this.onError) {
                this.onError({
                    code: 'RECONNECTING',
                    message: `Intentando reconectar (intento ${this.reconnectAttempts})`
                });
            }
            this.connect();
        }, delay);
    }

    subscribeToOrder(orderID) {
        if (!this.connected) {
            console.warn("Cannot subscribe: WebSocket not connected");
            return false;
        }

        this.orderID = orderID;
        const subscribeMessage = {
            type: 'SUBSCRIBE',
            order_id: orderID,
            timestamp: new Date()
        };

        console.log("Subscribing to order:", orderID);
        this.socket.send(JSON.stringify(subscribeMessage));
        return true;
    }

    disconnect() {
        console.log("Disconnecting WebSocket...");

        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
            this.reconnectTimeout = null;
        }

        if (this.socket) {
            this.socket.close(1000, 'Desconexión normal');
            this.socket = null;
        }

        this.connected = false;
    }
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = { DeliveryTracker, OrderTracker };
}
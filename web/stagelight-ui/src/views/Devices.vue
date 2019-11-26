<template>
    <table class="primary">
        <thead>
            <tr>
                <th>Device ID</th>
                <th>Color</th>
                <th>Mode</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="device in devices" v-bind:key="device.id">
                <td>{{ device.id }}</td>
                <td>{{ device.color }}</td>
                <td>{{ device.mode }}</td>
            </tr>
        </tbody>
    </table>
</template>

<script>
import axios from 'axios'

export default {
    name: 'Devices',
    data () {
        return {
            devices: []
        }
    },
    created () {
        axios.get('http://localhost:8000/device/all')
        .then(response => {
            this.devices = response.data
        })
    }
}
</script>
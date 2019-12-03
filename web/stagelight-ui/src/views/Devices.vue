<template>
    <div>
        <div class="select">
            <select v-model.lazy="selected.color">
                <option v-for="color in defaultColors" :key="color">{{ nameFromColor(color) }}</option>
            </select>
        </div>
        <div class="select">
            <select v-model="selected.mode">
                <option v-for="mode in modes" :key="mode">{{ nameFromModeID(mode) }}</option>
            </select>
        </div>
        <button class="button is-primary" v-on:click="applyDeviceChanges">Apply</button>
        <table class="table is-striped">
            <thead>
                <tr>
                    <th>Select</th>
                    <th>ID</th>
                    <th>Color</th>
                    <th>Mode</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="device in devices" :key="device.id">
                    <td>
                        <label class="checkbox">
                            <input type="checkbox" :value="device.id" v-model="selected.devices">
                        </label>
                    </td>
                    <td>{{ device.id }}</td>
                    <td v-bind:value="device.colors">{{ nameFromColor(device.colors[0]) }}</td>
                    <td v-bind:value="device.mode">{{ nameFromModeID(device.mode) }}</td>
                </tr>
            </tbody>
        </table>
        <button class="button is-primary" v-on:click="submitData">Submit</button>
    </div>
</template>

<script>
import axios from 'axios'

export default {
    name: 'Devices',
    data () {
        return {
            devices: [],
            defaultColors: {
                'Red': 16711680,
                'Green': 65280,
                'Blue': 255,
                'Off': 0
            },
            modes: {
                'Normal': 0,
                'Vote': 1
            },
            selected: {
                devices: [],
                color: 'Red',
                mode: 'Normal'
            }
        }
    },
    created () {
        axios.get(
            `${process.env.VUE_APP_API_BASE_URL}/device/all`
        ).then(response => {
            this.devices = response.data
        })
    },
    methods: {
        swap: function (obj) {
            var output = {}
            Object.assign(output, ...Object.entries(obj).map(([a, b]) => ({[b]: a})))
            return output
        },
        nameFromColor: function (color) {
            return this.swap(this.defaultColors)[color]
        },
        nameFromModeID: function (mode) {
            return this.swap(this.modes)[mode]
        },
        getOtherItems: function(obj, first) {
            var keys = Object.keys(obj).filter(e => (e != first))
            keys.unshift(first)

            return keys
        },
        submitData: function () {
            axios.post(
                `${process.env.VUE_APP_API_BASE_URL}/device/update`,
                this.devices
            )
        },
        applyDeviceChanges: function () {
            var i = 0
            this.devices.forEach(v => {
                if (this.selected.devices[i] === v.id) {
                    v.colors[0] = this.defaultColors[this.selected.color]
                    v.mode = this.modes[this.selected.mode]

                    i++
                }
            })
        }
    }
}
</script>
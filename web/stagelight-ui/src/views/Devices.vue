<template>
    <section class="section">
        <div class="container">
            <nav class="level">
                <div class="level-left">
                    <div class="level-item select">
                        <select v-model="selected.color">
                            <option v-for="color in colorSchemes" :key="color.id" :value="color.id">{{ color.name }}</option>
                        </select>
                    </div>
                    <div class="level-item select">
                        <select v-model="selected.mode">
                            <option v-for="mode in modes" :key="mode">{{ nameFromModeID(mode) }}</option>
                        </select>
                    </div>
                    <button class="button is-primary" v-on:click="applyDeviceChanges">Apply</button>
                </div>
                <div class="level-right">
                    <button class="button is-primary" v-on:click="submitData">Submit Changes</button>
                </div>
            </nav>
            <div class="box">
                <table class="table is-striped is-fullwidth">
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
                            <td v-bind:value="device.color">{{ getColorName(device.color) }}</td>
                            <td v-bind:value="device.mode">{{ nameFromModeID(device.mode) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </section>
</template>

<script>
import axios from 'axios'

export default {
    name: 'Devices',
    data () {
        return {
            devices: [],
            colorSchemes: [],
            modes: {
                'Normal': 0,
                'Vote': 1
            },
            selected: {
                devices: [],
                color: undefined,
                mode: undefined
            }
        }
    },
    created () {
        axios.get(
            `${process.env.VUE_APP_API_BASE_URL}/device/all`
        ).then(response => {
            this.devices = response.data
        })

        axios.get(
            `${process.env.VUE_APP_API_BASE_URL}/colors`
        ).then(response => {
            this.colorSchemes = response.data
        })
    },
    methods: {
        swap: function (obj) {
            var output = {}
            Object.assign(output, ...Object.entries(obj).map(([a, b]) => ({[b]: a})))
            return output
        },
        nameFromModeID: function (mode) {
            return this.swap(this.modes)[mode]
        },
        getColorName: function (id) {
            for (var i = 0; i < this.colorSchemes.length; i++) {
                if (this.colorSchemes[i].id === id) {
                    return this.colorSchemes[i].name
                }
            }

            return "Unknown"
        },
        submitData: function () {
            axios.post(
                `${process.env.VUE_APP_API_BASE_URL}/device/update`,
                this.devices
            )
        },
        applyDeviceChanges: function () {
            var i = 0
            var newColor = this.selected.color
            var newMode = this.modes[this.selected.mode]

            this.devices.forEach(v => {
                if (this.selected.devices[i] === v.id) {
                    if (typeof newColor !== "undefined") {
                        v.color = newColor
                    }

                    if (typeof newMode !== "undefined") {
                        v.mode = newMode
                    }

                    i++
                }
            })
        }
    }
}
</script>
<template>
    <div>
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
                    <td>
                        <select>
                            <option>{{ nameFromColor(device.colors[0]) }}</option>
                            <option v-for="opt in getOtherItems(defaultColors, nameFromColor(device.colors[0]))" v-bind:key="opt">
                                {{ opt }}
                            </option>
                        </select>
                    </td>
                    <td>
                        <select>
                            <option>{{ nameFromModeID(device.mode) }}</option>
                            <option v-for="opt in getOtherItems(modes, nameFromModeID(device.mode))" v-bind:key="opt">
                                {{ opt }}
                            </option>
                        </select>
                    </td>
                </tr>
            </tbody>
        </table>
        <button>Submit</button>
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
            }
        }
    },
    created () {
        axios.get('http://localhost:8000/device/all')
        .then(response => {
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
        getOtherItems: function(obj, rem) {
            var keys = Object.keys(obj)
            return keys.filter(e => (e != rem))
        }
    }
}
</script>
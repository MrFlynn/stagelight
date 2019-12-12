<template>
    <section class="section">
        <div class="container">
            <div class="columns">
                <div class="column is-half">
                    <h1 class="title">Color Sequences</h1>
                    <div class="tile is-ancestor">
                        <div class="tile is-vertical is-parent">
                            <div v-for="color in colors" class="tile is-child box" :key="color.id">
                                <h2 class="subtitle"><strong>Name:</strong> {{ color.name }}</h2>
                                <div class="content">
                                    <strong>Sequence:</strong>
                                    <ol>
                                        <div v-for="s in color.sequence" class="color-container" :key="s">
                                            <div>
                                                <div class="color-box" :style="{ backgroundColor: s.toStringPad()}"></div>
                                                <li>{{ s.toStringPad() }}</li>
                                            </div>
                                        </div>
                                    </ol>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="column is-half">
                    <h1 class="title">Add Color Sequence</h1>
                    <div class="box">
                        <div class="field">
                            <label class="label">Name</label>
                            <div class="control">
                                <input v-model="name" class="input" type="text" placeholder="Sequence name">
                            </div>
                        </div>
                        <div class="field">
                            <label class="label">Sequence</label>
                            <div class="control">
                                <textarea v-model="sequence" class="textarea" placeholder="Comma-separated RGB hex color sequences without leading #."></textarea>
                            </div>
                        </div>
                        <div class="field">
                            <div class="control">
                                <button class="button is-info" v-on:click="submit">Submit</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script>
import axios from 'axios'

export default {
    name: 'Colors',
    data () {
        return {
            colors: [],
            name: '',
            sequence: ''
        }
    },
    created () {
        axios.get(
            `${process.env.VUE_APP_API_BASE_URL}/colors`
        ).then(r => {
            this.colors = r.data
        })
    },
    methods: {
        submit: function () {
            const intSequence = this.sequence.split(',').map(e => parseInt(e, 16))
            var content = {
                name: this.name,
                sequence: intSequence
            }
            
            axios.post(
                `${process.env.VUE_APP_API_BASE_URL}/colors/new`,
                content
            )

            this.sequence = ''
            this.name = ''
        }
    }
}

Number.prototype.toStringPad = function () {
    var res = this.toString(16)
    return '#' + res.padStart(6, '0')
}

</script>

<style scoped>
    .color-container {
        position: relative;
    }

    .color-container .color-box {
        float: right;
        width: 20px;
        height: 20px;
        background-color: black;
        clear: both;
    }
</style>
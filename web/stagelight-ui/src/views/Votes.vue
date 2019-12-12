<template>
    <section class="section">
        <div class="container">
            <div class="columns">
                <div class="column is-half">
                    <h1 class="title">Vote Progress Bar</h1>
                    <progress class="progress is-primary" :value="percent" max="100"></progress>
                </div>
                <div class="column is-half">
                    <h1 class="title">Vote Counts</h1>
                    <div class="tile is-ancestor">
                        <div class="tile is-parent">
                            <div class="tile is-child notification is-success">
                                <div class="content">
                                    <h2 class="subtitle">Positive votes</h2>
                                    <p class="subtitle is-5">{{ votes["positive"] }}</p>
                                </div>
                            </div>
                            <div class="tile is-child notification is-danger">
                                <div class="content">
                                    <h2 class="subtitle">Negative votes</h2>
                                    <p class="subtitle is-5">{{ votes["negative"] }}</p>
                                </div>
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
    name: 'Votes',
    data () {
        return {
            votes: {},
            percent: 0
        }
    },
    created () {
        axios.get(
            `${process.env.VUE_APP_API_BASE_URL}/votes`
        ).then(r => {
            this.votes = r.data,
            this.percent = Math.floor((this.votes["positive"] / (this.votes["positive"] + this.votes["negative"])) * 100)
        })
    }
}
</script>
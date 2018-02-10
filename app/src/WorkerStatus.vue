<template>
    <div class="col-xs-12 col-md-6 col-lg-4">
        <div class="box box-solid" :class="boxClassName">
            <div class="box-header">
                <i class="fa" :class="icon"></i>
                <h3 class="box-title">{{type}}</h3>
                <span class="pull-right badge bg-orange">{{workerCount}} online</span>
            </div>
            <div class="box-body">
                <table class="table table-hover">
                    <tbody>
                    <tr v-for="(worker, index) in workers">
                        <td>
                            <i class="fa fa-heartbeat li-icon"></i><span class="li-text">{{worker.Name}}</span>
                            <span class="badge-center pull-right badge time-text">{{formatDateTime(worker.Time * 1000)}}</span>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script>
    import formatDateTime from "./datetimeUtil.js"
    export default {
        name: "worker-status",
        props: [
            "type",
            "icon",
            "workers"
        ],
        data() {
            return {}
        },
        computed: {
            workerCount() {
                let vm = this;
                let c = 0;
                for (let w in vm.workers)
                    c++;
                return c;
            },
            boxClassName() {
                let vm = this;
                if (vm.workerCount <= 0)
                    return "box-danger";
                else
                    return "box-success";
            }
        },
        methods: {
            formatDateTime(time) {
                return formatDateTime(time);
            }
        }
    }
</script>

<style scoped>
    .li-text {
        padding-left: 0.5em;
    }

    .li-icon {
        padding-top: 0.8em;
    }

    .badge-center {
        margin-top: 0.6em;
    }
</style>
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
                            <span class="badge-center pull-right badge">{{timeStamp2String(worker.Time)}}</span>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script>
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
            timeStamp2String(time) {
                let datetime = new Date(time * 1000);
                let year = datetime.getFullYear();
                let month = datetime.getMonth() + 1;
                let date = datetime.getDate();
                let hour = datetime.getHours();
                let minute = datetime.getMinutes();
                let second = datetime.getSeconds();
                return year + "-" + month + "-" + date + "T" + hour + ":" + minute + ":" + second + "Z";
            },
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

    .status-text {
        margin-left: 0.3em;
        font-size: 0.4em;
    }

    .btn-center {
        margin-top: -1.5em;
    }

    .badge-center {
        margin-top: 0.6em;
    }
</style>
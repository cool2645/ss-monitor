<template>
    <div class="col-xs-12 col-md-6 col-lg-4">
        <div class="box box-solid" :class="boxClassName">
            <div class="box-header">
                <i class="fa fa-cube"></i>
                <h3 class="box-title">{{name}}</h3>
                <span class="status-text">
                    <span>丢包率 </span>
                    <span>{{avgLossTime}}</span>
                </span>
                <span class="status-text">
                    <span>延迟 </span>
                    <span>{{avgResTime}}</span>
                </span>
                <span v-show="isCleaning" class="pull-right badge bg-orange">正在清洗</span>
            </div>
            <div class="box-body">
                <table class="table table-hover">
                    <tbody>
                    <tr>
                        <td>
                            <span>
                                <i class="fa fa-sitemap li-icon"></i><span class="li-text">CN 可访问性</span>
                                <br/><span class="status-text">{{cnTaskTime}}</span>
                            </span>
                            <a class="btn-center pull-right btn"
                               :class="getTaskButtonClassName(cnTaskState)"
                               @click="location(getTaskLink(cnTaskID))">{{cnTaskState}}</a>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <span>
                                <i class="fa fa-send li-icon"></i><span class="li-text">SS (IPv4) 可用性</span>
                                <br/><span class="status-text">{{ssTaskTime}}</span>
                            </span>
                            <a class="btn-center pull-right btn"
                               :class="getTaskButtonClassName(ssTaskState)"
                               @click="location(getTaskLink(ssTaskID))">{{ssTaskState}}</a>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <span>
                                <i class="fa fa-send li-icon"></i><span class="li-text">SS (IPv6) 可用性</span>
                                <br/><span class="status-text">{{ss6TaskTime}}</span>
                            </span>
                            <a class="btn-center pull-right btn"
                               :class="getTaskButtonClassName(ss6TaskState)"
                               @click="location(getTaskLink(ss6TaskID))">{{ss6TaskState}}</a>
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
        name: "node-status",
        props: [
            'node'
        ],
        data() {
            return {
                name: this.node.Name,
                avgLossTime: eval('(' + this.node.Status.CN.Result + ')').avg_lost_percent.toFixed(2),
                avgResTime: eval('(' + this.node.Status.CN.Result + ')').avg_res_time.toFixed(2),
                isCleaning: this.node.IsCleaning,
                cnTaskID: this.node.Status.CN.ID,
                ssTaskID: this.node.Status.SS.ID,
                ss6TaskID: this.node.Status['SS-IPv6'].ID,
            }
        },
        computed: {
            boxClassName() {
                let c = "box-success";
                let vm = this;
                if (vm.isCleaning || vm.cnTaskState === "Failing"
                    || vm.ssTaskState === "Failing" || vm.ss6TaskState === "Failing") {
                    c = "box-danger";
                }
                return c;
            },
            cnTaskState() {
                let vm = this;
                if (vm.cnTaskID === 0)
                    return "Unknown";
                else
                    return vm.node.Status.CN.State;
            },
            cnTaskTime() {
                let vm = this;
                if (vm.cnTaskID === 0)
                    return "Never";
                else
                    return vm.node.Status.CN.UpdatedAt;
            },
            ssTaskState() {
                let vm = this;
                if (vm.ssTaskID === 0)
                    return "Unknown";
                else
                    return vm.node.Status.SS.State;
            },
            ssTaskTime() {
                let vm = this;
                if (vm.ssTaskID === 0)
                    return "Never";
                else
                    return vm.node.Status.SS.UpdatedAt;
            },
            ss6TaskState() {
                let vm = this;
                if (vm.ss6TaskID === 0)
                    return "Unknown";
                else
                    return vm.node.Status['SS-IPv6'].State;
            },
            ss6TaskTime() {
                let vm = this;
                if (vm.ss6TaskID === 0)
                    return "Never";
                else
                    return vm.node.Status['SS-IPv6'].UpdatedAt;
            },
        },
        methods: {
            location(href) {
                this.$router.push({path: href});
            },
            getTaskButtonClassName(state) {
                let c = "btn-warning";
                if (state === "Passing")
                    c = "btn-success";
                else if (state === "Failing")
                    c = "btn-danger";
                else if (state === "Queuing")
                    c = "btn-info";
                else if (state === "Shiny☆")
                    c = "btn-shiny";
                else if (state === "Unknown")
                    c = "btn-default";
                return c;
            },
            getTaskLink(id) {
                if (id === 0)
                    return '/launch';
                else
                    return '/task/' + id;
            }
        },
        mounted() {
            // console.log(this.ss6TaskID);
        }
    }
</script>

<style scoped lang="scss">
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

    .btn-shiny {
        color: #fff;
        background-color: #b76add;
        border-color: #a761ca;
        &:hover {
            color: #fff;
            background-color: #a761ca;
            border-color: #a761ca;
        }
    }
</style>
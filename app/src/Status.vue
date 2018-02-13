<template>
    <div>
        <!-- Content Header (Page header) -->
        <section class="content-header">
            <h1>
                状态监控
                <small>Real-time Status</small>
            </h1>
        </section>
        <!-- Content -->
        <section class="content">
            <!--Error Message-->
            <div class="row">
                <div class="col-md-12">
                    <div id="msg-success" class="alert alert-success alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-info"></i> 成功!</h4>

                        <p id="msg-success-p"></p>
                    </div>
                    <div id="msg-warning" class="alert alert-warning alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-warning-p">{{ warning }}</p>
                    </div>
                    <div id="msg-error" class="alert alert-danger alert-dismissable" style="display: none;">
                        <button type="button" class="close" @click="dismissAlert" aria-hidden="true">&times;</button>
                        <h4><i class="icon fa fa-warning"></i> 出错了!</h4>

                        <p id="msg-error-p">电波无法送达哦～<a href="javascript:;" @click="startClocking">重试</a></p>
                    </div>
                </div>
            </div>
            <!--./Error Message-->
            <!--Nodes-->
            <section class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-cubes"></i>
                            <h3 class="box-title">节点状态</h3>
                        </div>
                        <div class="box-body">
                            <!--Node status-->
                            <transition v-for="(node, index) in nodes" :key="index" name="fade" v-show="loadFinish">
                                <node-status :node="node"></node-status>
                            </transition>
                            <!--./Node status-->
                        </div>
                    </div>
                </div>
            </section>
            <!--./Nodes-->
            <!--Workers-->
            <section class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-group"></i>
                            <h3 class="box-title">Worker 状态</h3>
                        </div>
                        <div class="box-body">
                            <!--Watcher status-->
                            <transition name="fade">
                                <worker-status type="Watcher" icon="fa-eye" v-show="loadFinish"
                                               :workers="workers.watcher"></worker-status>
                            </transition>
                            <!--./Watcher status-->
                            <!--IPv4 Tester status-->
                            <transition name="fade">
                                <worker-status type="IPv4 Tester" icon="fa-gear" v-show="loadFinish"
                                               :workers="workers.ipv4tester"></worker-status>
                            </transition>
                            <!--./IPv4 Tester status-->
                            <!--IPv6 Tester status-->
                            <transition name="fade">
                                <worker-status type="IPv6 Tester" icon="fa-gear" v-show="loadFinish"
                                               :workers="workers.ipv6tester"></worker-status>
                            </transition>
                            <!--./IPv6 Tester status-->
                            <!--Cleaner status-->
                            <transition name="fade">
                                <worker-status type="Cleaner" icon="fa-ambulance" v-show="loadFinish"
                                               :workers="workers.cleaner"></worker-status>
                                <!--./Cleaner status-->
                            </transition>
                        </div>
                    </div>
                </div>
            </section>
            <!--./Workers-->
            <div class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-list"></i>
                            <h3 class="box-title">详细信息</h3>
                        </div>
                        <div class="box-body">
                            <tree-view :data="jsonSource" :options="{maxDepth: 0}"></tree-view>
                        </div>
                    </div>
                </div>
            </div>
        </section>
        <!-- /.content -->
    </div><!-- /.content-wrapper -->
</template>

<script>
    import config from './config'
    import NodeStatus from "./NodeStatus.vue"
    import WorkerStatus from "./WorkerStatus.vue";

    export default {
        components: {
            WorkerStatus,
            NodeStatus
        },
        data() {
            return {
                jsonSource: {},
                nodes: {},
                workers: {},
                loadFinish: false,
                warning: "",
            }
        },
        computed: {
            //
        },
        component: {
            'node-status': NodeStatus,
            'worker-status': WorkerStatus,
        },
        created() {
            let vm = this;
            vm.jsonSource = {};
        },
        mounted() {
            this.startClocking();
        },
        methods: {
            dismissAlert() {
                $("#msg-success").hide(10);
                $("#msg-warning").hide(10);
                $("#msg-error").hide(10);
            },
            startClocking() {
                $("#msg-error").hide(10);
                this.updateData(true);
            },
            updateData(recur) {
                let vm = this;
                fetch(config.urlPrefix + '/status')
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSource = res;
                                    vm.nodes = vm.jsonSource.data.nodes;
                                    vm.workers = vm.jsonSource.data.workers;
                                    vm.loadFinish = true;
                                    if(!vm._isBeingDestroyed && recur) setTimeout(() => {this.updateData(true)}, 5000);
                                    vm.$emit('async-load');
                                }
                            }
                        )
                    })
                    .catch(error => {
                        $("#msg-error").hide(10).show(100);
                    });
            },
        }
    }
</script>

<style scoped>
    .fade-enter-active, .fade-leave-active {
        transition: opacity .5s;
    }

    .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */
    {
        opacity: 0;
    }
</style>
<template>
    <div>
        <section class="content-header">
            <h1>
                创建任务
                <small>Launch Task</small>
            </h1>
        </section>

        <!-- Main content -->
        <section class="content">
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

                        <p id="msg-error-p">{{ error }}</p>
                    </div>
                </div>
            </div>
            <div class="row">
                <!-- left column -->
                <div class="col-md-6">
                    <!-- general form elements -->
                    <div class="box box-primary">
                        <div class="box-body">
                            <div class="form-horizontal">
                                <div class="row">
                                    <fieldset class="col-sm-12">
                                        <legend>创建 Watcher 任务</legend>
                                        <p>Watcher 任务通过 Ping 来测试服务器的中国大陆可访问性。</p>
                                        <p>请在服务器名中填入域名或 IPv4 地址。</p>
                                        <div class="form-group">
                                            <label for="watcher-server-name" class="col-sm-3 control-label">服务器名</label>

                                            <div class="col-sm-9">
                                                <input class="form-control" id="watcher-server-name" v-model="watcher.server_name">
                                            </div>
                                        </div>

                                    </fieldset>
                                </div>
                            </div>
                        </div>
                        <!-- /.box-body -->
                        <div class="box-footer">
                            <button @click="launchWatcher" class="btn btn-primary">添加</button>
                        </div>
                    </div>
                    <!-- /.box -->
                </div>
                <!-- right column -->
                <div class="col-md-6">
                    <!-- general form elements -->
                    <div class="box box-primary">
                        <div class="box-body">
                            <div class="form-horizontal">
                                <div class="row">
                                    <fieldset class="col-sm-12">
                                        <legend>创建 Tester 任务</legend>
                                        <p>Tester 任务通过 Docker 下的 Shadowsocks 客户端来测试 Shadowsocks 服务器的响应情况。</p>
                                        <p>请在 Json 配置一栏中填入 Shadowsocks 的客户端配置。</p>
                                        <p>请在服务器名中填入友好的对该服务器的称呼。你的 Json 配置不会被公开，我们将用你填写的“服务器名”来称呼它。</p>

                                        <div class="form-group">
                                            <label for="tester-server-name" class="col-sm-3 control-label">服务器名</label>

                                            <div class="col-sm-9">
                                                <input class="form-control" id="tester-server-name" v-model="tester.server_name">
                                            </div>
                                        </div>

                                        <div class="form-group">
                                            <label for="ip-ver" class="col-sm-3 control-label">IP 协议</label>

                                            <div class="col-sm-9">
                                                <select class="form-control" id="ip-ver" v-model="tester.ip_ver">
                                                    <option value="4">IPv4</option>
                                                    <option value="6">IPv6</option>
                                                </select>
                                            </div>
                                        </div>

                                        <div class="form-group">
                                            <label for="ss-json" class="col-sm-3 control-label">Json 配置</label>

                                            <div class="col-sm-9">
                                                <textarea class="form-control" id="ss-json" rows="6" v-model="tester.ss_json"></textarea>
                                            </div>
                                        </div>
                                    </fieldset>
                                </div>
                            </div>
                        </div>
                        <!-- /.box-body -->
                        <div class="box-footer">
                            <button @click="launchTester" class="btn btn-primary">添加</button>
                        </div>
                    </div>
                    <!-- /.box -->
                </div>
                <!-- /.row -->
            </div>
            <!-- /.row -->
        </section>
        <!-- /.content -->
    </div>
</template>

<script>
    import config from './config'
    import urlParam from './buildUrlParam'
    export default {
        data() {
            return {
                warning: "",
                error: "",
                watcher: {
                    class: "watcher",
                    server_name: ""
                },
                tester: {
                    class: "tester",
                    ip_ver: 4,
                    server_name: "",
                    ss_json: `{
    "server": "93.184.216.34",
    "server_port": 8388,
    "password": "shadowsocks",
    "method": "aes-128-gcm"
}`
                }
            }
        },
        methods: {
            dismissAlert() {
                $("#msg-success").hide(10);
                $("#msg-warning").hide(10);
                $("#msg-error").hide(10);
            },
            launchWatcher() {
                if (!this.watcher.server_name) {
                    this.warning = "请填写服务器名";
                    $("#msg-warning").hide(10).show(100);
                    return
                }
                this.launchTask(this.watcher)
            },
            launchTester() {
                if (!this.tester.server_name) {
                    this.warning = "请填写服务器名";
                    $("#msg-warning").hide(10).show(100);
                    return
                } else if(!this.tester.ss_json) {
                    this.warning = "请填写 Json 配置";
                    $("#msg-warning").hide(10).show(100);
                    return
                }
                this.launchTask(this.tester)
            },
            launchTask(worker) {
                let vm = this;
                fetch(config.urlPrefix + '/task?', {
                    method: "POST",
                    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                    body: urlParam(worker)
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    $("#msg-warning").hide(100);
                                    $("#msg-error").hide(100);
                                    $("#msg-success").show(100);
                                    window.setTimeout("location.href = '#/task'", 2000);
                                } else {
                                    vm.error = "发生错误：" + res.msg;
                                    $("#msg-warning").hide(100);
                                    $("#msg-error").hide(10).show(100);
                                }
                            }
                        )
                    });
            }
        }
    }
</script>

<style>

</style>
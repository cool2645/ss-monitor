<template>
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
                                    <legend>管理员登录</legend>
                                    <div class="form-group">
                                        <label for="username" class="col-sm-3 control-label">用户名</label>

                                        <div class="col-sm-9">
                                            <input class="form-control" id="username">
                                        </div>
                                    </div>
                                    <div class="form-group">
                                        <label for="password" class="col-sm-3 control-label">密码</label>

                                        <div class="col-sm-9">
                                            <input class="form-control" id="password">
                                        </div>
                                    </div>
                                </fieldset>
                            </div>
                        </div>
                    </div>
                    <!-- /.box-body -->
                    <div class="box-footer">
                        <button @click="" class="btn btn-primary">登录</button>
                    </div>
                </div>
                <!-- /.box -->
            </div>
        </div>
        <!-- /.row -->
    </section>
    <!-- /.content -->
</template>

<script>
    export default {
        data() {
            return {
                warning: "",
                error: "",
            }
        },
        methods: {
            dismissAlert() {
                $("#msg-success").hide(10);
                $("#msg-warning").hide(10);
                $("#msg-error").hide(10);
            },
        }
    }
</script>

<style scoped>

</style>
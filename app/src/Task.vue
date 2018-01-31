<template>
    <div>
        <!-- Content Header (Page header) -->
        <section class="content-header">
            <h1>
                任务状态
                <small>Task Status</small>
            </h1>
        </section>
        <!-- Content -->
        <section class="content">
            <div class="row">
                <div class="col-xs-12">
                    <div class="box box-default">
                        <div class="box-header">
                            <i class="fa fa-list"></i>
                            <h3 class="box-title">详细信息</h3>
                        </div>
                        <div class="box-body">
                            <tree-view :data="jsonSource" :options="{maxDepth: 3}"></tree-view>
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
    export default {
        data() {
            return {
                jsonSource: {},
            }
        },
        mounted() {
            this.updateData()
        },
        methods: {
            updateData() {
                let vm = this;
                fetch(config.urlPrefix + '/task/' + this.$route.params.id)
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.jsonSource = res
                                }
                            }
                        )
                    });
            }
        },
        watch: {
            '$route' (to, from) {
                this.updateData()
            }
        }
    }
</script>

<style scoped>

</style>
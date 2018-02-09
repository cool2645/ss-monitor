<template>
    <div class="col-md-12">
        <div :class="boxClass">
            <div class="box-body">
                <ul class="products-list product-list-in-box">
                    <li class="item">
                        <a v-if="!editing" href="javascript:;" class="product-title">{{ nodeDisplayName }}</a>
                        <input v-else v-model="data.Name" class="product-title"/>
                        <a v-if="!editing" @click="edit" class="label label-info pull-right">编辑</a>
                        <a v-else @click="update" class="label label-info pull-right">保存</a>
                        <a v-if="!editing" @click="destroy" class="label label-danger pull-right">删除</a>
                        <a v-else @click="quit" class="label label-danger pull-right">放弃</a>
                        <a v-if="!isNew && node1.IsCleaning" @click="reset" class="label label-primary pull-right">重置</a>
                        <a v-else-if="!isNew && !node1.IsCleaning" @click="setCleaning" class="label label-warning pull-right">置清洗</a>
                    </li><!-- /.item -->
                </ul>
            </div>
            <div class="box-footer no-padding">
                <div class="row">
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>IPv4 域名前缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainPrefix4 }}</span>
                                <input v-else v-model="data.DomainPrefix4" class="pull-right"/>
                            </li>
                            <li>IPv6 域名前缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainPrefix6 }}</span>
                                <input v-else v-model="data.DomainPrefix6" class="pull-right"/>
                            </li>
                            <li>域名后缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainRoot }}</span>
                                <input v-else v-model="data.DomainRoot" class="pull-right"/>
                            </li>
                            <li>IPv4 地址
                                <span v-if="!editing" class="pull-right badge bg-green">{{ data.IPv4 }}</span>
                                <input v-else v-model="data.IPv4" class="pull-right"/>
                            </li>
                            <li>IPv6 地址
                                <span v-if="!editing" class="pull-right badge bg-green">{{ data.IPv6 }}</span>
                                <input v-else v-model="data.IPv6" class="pull-right"/>
                            </li>
                            <li>更新时间 <span class="pull-right badge bg-grey">{{ data.UpdatedAt }}</span></li>
                        </ul>
                    </div>
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>云计算服务商
                                <span v-if="!editing" class="pull-right badge bg-maroon">{{ data.Provider }}</span>
                                <input v-else v-model="data.Provider" class="pull-right"/>
                            </li>
                            <li>DNS 服务商
                                <span v-if="!editing" class="pull-right badge bg-maroon">{{ data.DNSProvider }}</span>
                                <input v-else v-model="data.DNSProvider" class="pull-right"/>
                            </li>
                            <li>数据中心
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DataCenter }}</span>
                                <input v-else v-model="data.DataCenter" class="pull-right"/>
                            </li>
                            <li>规格
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.Plan }}</span>
                                <input v-else v-model="data.Plan" class="pull-right"/>
                            </li>
                            <li>操作系统
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.OS }}</span>
                                <input v-else v-model="data.OS" class="pull-right"/>
                            </li>
                            <li>镜像
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.Image }}</span>
                                <input v-else v-model="data.Image" class="pull-right"/>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv4 Json
                                <p v-if="!editing" class="pull-right">
                                    {{ data.Ss4Json }}
                                </p>
                                <textarea v-else v-model="data.Ss4Json" class="pull-right"></textarea>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv6 Json
                                <p v-if="!editing" class="pull-right">
                                    {{ data.Ss6Json }}
                                </p>
                                <textarea v-else v-model="data.Ss6Json" class="pull-right"></textarea>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 Watch
                                <span v-if="!editing" class="pull-right badge bg-aqua">{{ data.EnableWatching }}</span>
                                <input v-else v-model="data.EnableWatching" type="checkbox" class="pull-right"/>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv4 Test
                                <span v-if="!editing"
                                      class="pull-right badge bg-aqua">{{ data.EnableIPv4Testing }}</span>
                                <input v-else v-model="data.EnableIPv4Testing" type="checkbox" class="pull-right"/>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv6 Test
                                <span v-if="!editing"
                                      class="pull-right badge bg-aqua">{{ data.EnableIPv6Testing }}</span>
                                <input v-else v-model="data.EnableIPv6Testing" type="checkbox" class="pull-right"/>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许清洗
                                <span v-if="!editing" class="pull-right badge bg-aqua">{{ data.EnableCleaning }}</span>
                                <input v-else v-model="data.EnableCleaning" type="checkbox" class="pull-right"/>
                            </li>
                        </ul>
                    </div>
                </div>

            </div>
        </div><!-- /.widget-user -->
    </div><!-- /.col -->
</template>

<script>
    import config from './config'

    export default {
        props: [
            "node",
            "new"
        ],
        data() {
            return {
                editing: false,
                data: {},
                node1: this.node,
                isNew: this.new
            }
        },
        mounted() {
            this.data = Object.assign({}, this.node1);
            if (this.isNew)
                this.editing = true;
        },
        computed: {
            nodeDisplayName() {
                return this.node1.IsCleaning ? this.data.Name + "（清洗中）" : this.data.Name
            },
            boxClass() {
                return this.node1.IsCleaning ? "box box-warning" : "box box-primary"
            }
        },
        methods: {
            edit() {
                this.editing = true;
            },
            quit() {
                if (this.isNew) {
                    this.$emit('delete-node', this.key)
                }
                this.editing = false;
                this.data = Object.assign({}, this.node1);
            },
            update() {
                let vm = this;
                if (this.isNew)
                    fetch(config.urlPrefix + '/node', {
                        credentials: 'include',
                        method: "POST",
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify(this.data)
                    })
                        .then(res => {
                            res.json().then(
                                res => {
                                    if (res.result) {
                                        this.editing = false;
                                        vm.node1 = res.data;
                                        vm.data = Object.assign({}, vm.node1);
                                        this.isNew = false;
                                    } else if (res.code === 401) {
                                        this.$emit('show-warning', "登录超时：" + res.msg);
                                        this.$emit('check-auth');
                                    } else {
                                        this.$emit('show-warning', "发生错误：" + res.msg);
                                    }
                                }
                            )
                        })
                        .catch(error => {
                            this.$emit('show-error', "发生错误：" + res.msg);
                        });
                else
                    fetch(config.urlPrefix + '/node/' + this.data.ID, {
                        credentials: 'include',
                        method: "PUT",
                        headers: {'Content-Type': 'application/json'},
                        body: JSON.stringify(this.data)
                    })
                        .then(res => {
                            res.json().then(
                                res => {
                                    if (res.result) {
                                        this.editing = false;
                                        vm.node1 = res.data;
                                        vm.data = Object.assign({}, vm.node1);
                                    } else if (res.code === 401) {
                                        this.$emit('show-warning', "登录超时：" + res.msg);
                                        this.$emit('check-auth');
                                    } else {
                                        this.$emit('show-warning', "发生错误：" + res.msg);
                                    }
                                }
                            )
                        })
                        .catch(error => {
                            this.$emit('show-error', "发生错误：" + res.msg);
                        });
            },
            destroy() {
                let r = confirm("确定要删除节点吗？");
                if (!r) return;
                fetch(config.urlPrefix + '/node/' + this.data.ID, {
                    credentials: 'include',
                    method: "DELETE",
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'}
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    this.$emit('delete-node', this.key)
                                } else if (res.code === 401) {
                                    this.$emit('show-warning', "登录超时：" + res.msg);
                                    this.$emit('check-auth');
                                } else {
                                    this.$emit('show-warning', "发生错误：" + res.msg);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        this.$emit('show-error', "发生错误：" + res.msg);
                    });
            },
            reset() {
                let r = confirm("确定要重置清洗状态吗？");
                if (!r) return;
                let vm = this;
                fetch(config.urlPrefix + '/node/' + this.data.ID + "/status/isCleaning", {
                    credentials: 'include',
                    method: "DELETE",
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'}
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.node1.IsCleaning = false;
                                    vm.data = Object.assign({}, vm.node1);
                                } else if (res.code === 401) {
                                    this.$emit('show-warning', "登录超时：" + res.msg);
                                    this.$emit('check-auth');
                                } else {
                                    this.$emit('show-warning', "发生错误：" + res.msg);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        this.$emit('show-error', "发生错误：" + res.msg);
                    });
            },
            setCleaning() {
                let r = confirm("确定要置清洗状态吗？");
                if (!r) return;
                let vm = this;
                fetch(config.urlPrefix + '/node/' + this.data.ID + "/status/isCleaning", {
                    credentials: 'include',
                    method: "POST",
                    headers: {'Content-Type': 'application/x-www-form-urlencoded'}
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.node1.IsCleaning = true;
                                    vm.data = Object.assign({}, vm.node1);
                                } else if (res.code === 401) {
                                    this.$emit('show-warning', "登录超时：" + res.msg);
                                    this.$emit('check-auth');
                                } else {
                                    this.$emit('show-warning', "发生错误：" + res.msg);
                                }
                            }
                        )
                    })
                    .catch(error => {
                        this.$emit('show-error', "发生错误：" + res.msg);
                    });
            }
        }
    }
</script>

<style scoped>
    li {
        padding: 10px;

    }

    .pull-right {
        max-width: 60%;
    }

    textarea.pull-right {
        width: 70%;
        max-width: 70%;
    }

    span {
        word-break: break-all;
        display: block;
        white-space: normal;
    }
</style>
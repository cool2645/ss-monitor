<template>
    <div class="col-md-12">
        <div class="box box-primary">
            <div class="box-body">
                <ul class="products-list product-list-in-box">
                    <li class="item">
                        <a href="javascript:;" class="product-title">{{ data.Name }}</a>
                        <a v-if="!editing" @click="edit" class="label label-info pull-right">编辑</a>
                        <a v-else @click="update" class="label label-info pull-right">保存</a>
                        <a v-if="!editing" class="label label-danger pull-right">删除</a>
                        <a v-else @click="quit" class="label label-danger pull-right">放弃</a>
                    </li><!-- /.item -->
                </ul>
            </div>
            <div class="box-footer no-padding">
                <div class="row">
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>IPv4 域名前缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainPrefix4 }}</span>
                                <input v-else v-model="data.DomainPrefix4" class="pull-right" />
                            </li>
                            <li>IPv6 域名前缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainPrefix6 }}</span>
                                <input v-else v-model="data.DomainPrefix6" class="pull-right" />
                            </li>
                            <li>域名后缀
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DomainRoot }}</span>
                                <input v-else v-model="data.DomainRoot" class="pull-right" />
                            </li>
                            <li>IPv4 地址
                                <span v-if="!editing" class="pull-right badge bg-green">{{ data.IPv4 }}</span>
                                <input v-else v-model="data.IPv4" class="pull-right" />
                            </li>
                            <li>IPv6 地址
                                <span v-if="!editing" class="pull-right badge bg-green">{{ data.IPv6 }}</span>
                                <input v-else v-model="data.IPv6" class="pull-right" />
                            </li>
                            <li>更新时间 <span class="pull-right badge bg-grey">{{ data.UpdatedAt }}</span></li>
                        </ul>
                    </div>
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>云计算服务商
                                <span v-if="!editing" class="pull-right badge bg-maroon">{{ data.Provider }}</span>
                                <input v-else v-model="data.Provider" class="pull-right" />
                            </li>
                            <li>DNS 服务商
                                <span v-if="!editing" class="pull-right badge bg-maroon">{{ data.DNSProvider }}</span>
                                <input v-else v-model="data.DNSProvider" class="pull-right" />
                            </li>
                            <li>数据中心
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.DataCenter }}</span>
                                <input v-else v-model="data.DataCenter" class="pull-right" />
                            </li>
                            <li>规格
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.Plan }}</span>
                                <input v-else v-model="data.Plan" class="pull-right" />
                            </li>
                            <li>操作系统
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.OS }}</span>
                                <input v-else v-model="data.OS" class="pull-right" />
                            </li>
                            <li>镜像
                                <span v-if="!editing" class="pull-right badge bg-blue">{{ data.Image }}</span>
                                <input v-else v-model="data.Image" class="pull-right" />
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv4 Json
                                <p  v-if="!editing"  class="pull-right">
                                    {{ data.Ss4Json }}
                                </p>
                                <textarea v-else v-model="data.Ss4Json" class="pull-right"></textarea>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv4 Json
                                <p v-if="!editing"  class="pull-right">
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
                                <input v-else v-model="data.EnableWatching" type="checkbox" class="pull-right" />
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv4 Test
                                <span v-if="!editing" class="pull-right badge bg-aqua">{{ data.EnableIPv4Testing }}</span>
                                <input v-else v-model="data.EnableIPv4Testing" type="checkbox" class="pull-right" />
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv6 Test
                                <span v-if="!editing" class="pull-right badge bg-aqua">{{ data.EnableIPv6Testing }}</span>
                                <input v-else v-model="data.EnableIPv6Testing" type="checkbox" class="pull-right" />
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许清洗
                                <span v-if="!editing" class="pull-right badge bg-aqua">{{ data.EnableCleaning }}</span>
                                <input v-else v-model="data.EnableCleaning" type="checkbox" class="pull-right" />
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
            "node"
        ],
        data() {
            return {
                editing: false,
                data: {}
            }
        },
        mounted() {
            this.data = Object.assign({}, this.node);
        },
        methods: {
            edit() {
                this.editing = true;
            },
            quit() {
                this.editing = false;
                this.data = Object.assign({}, this.node);
            },
            update() {
                let vm = this;
                fetch(config.urlPrefix + '/node/' + this.data.ID,  {
                    credentials: 'include',
                    method: "PUT",
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(this.data)
                })
                    .then(res => {
                        res.json().then(
                            res => {
                                if (res.result) {
                                    vm.node = res.data;
                                    vm.data = Object.assign({}, vm.node);
                                }
                            }
                        )
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
<template>
    <div class="col-md-12">
        <div class="box box-primary">
            <div class="box-body">
                <ul class="products-list product-list-in-box">
                    <li class="item">
                        <a href="javascript:;" class="product-title">节点名</a>
                        <a v-if="editing" class="label label-info pull-right">编辑</a>
                        <a v-else class="label label-info pull-right">保存</a>
                    </li><!-- /.item -->
                </ul>
            </div>
            <div class="box-footer no-padding">
                <div class="row">
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>IPv4 域名前缀
                                <span v-if="editing" class="pull-right badge bg-blue">fake.us</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>IPv6 域名前缀
                                <span v-if="editing" class="pull-right badge bg-blue">fake.us6</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>域名后缀
                                <span v-if="editing" class="pull-right badge bg-blue">example.domain</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>IPv4 地址
                                <span v-if="editing" class="pull-right badge bg-green">208.167.255.14</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>IPv6 地址
                                <span v-if="editing" class="pull-right badge bg-green">2001:19f0:5:2b9c:5400:01ff:fe5c:3af1</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>更新时间 <span class="pull-right badge bg-grey">2018-02-04T11:41:09Z</span></li>
                        </ul>
                    </div>
                    <div class="col-sm-6">
                        <ul class="nav nav-stacked">
                            <li>云计算服务商
                                <span v-if="editing" class="pull-right badge bg-maroon">Vultr</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>DNS 服务商
                                <span v-if="editing" class="pull-right badge bg-maroon">DNSimple</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>数据中心
                                <span v-if="editing" class="pull-right badge bg-blue">1</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>规格
                                <span v-if="editing" class="pull-right badge bg-blue">201</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>操作系统
                                <span v-if="editing" class="pull-right badge bg-blue">FromSnapshot</span>
                                <input v-else class="pull-right" />
                            </li>
                            <li>镜像
                                <span v-if="editing" class="pull-right badge bg-blue">123wedvfdfds</span>
                                <input v-else class="pull-right" />
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv4 Json
                                <p  v-if="editing"  class="pull-right">
                                    {
                                    "server": "93.184.216.34",
                                    "server_port": 8388,
                                    "password": "shadowsocks",
                                    "method": "aes-128-gcm"
                                    }
                                </p>
                                <textarea v-else class="pull-right"></textarea>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-12">
                        <ul class="nav nav-stacked">
                            <li>
                                SS IPv4 Json
                                <p v-if="editing"  class="pull-right">
                                    {
                                    "server": "93.184.216.34",
                                    "server_port": 8388,
                                    "password": "shadowsocks",
                                    "method": "aes-128-gcm"
                                    }
                                </p>
                                <textarea v-else class="pull-right"></textarea>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 Watch
                                <span class="pull-right badge bg-aqua">1</span>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv4 Test
                                <span class="pull-right badge bg-aqua">1</span>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许 IPv6 Test
                                <span class="pull-right badge bg-aqua">1</span>
                            </li>
                        </ul>
                    </div>
                    <div class="col-sm-6 col-md-3">
                        <ul class="nav nav-stacked">
                            <li>允许清洗
                                <span class="pull-right badge bg-aqua">1</span>
                            </li>
                        </ul>
                    </div>
                </div>

            </div>
        </div><!-- /.widget-user -->
    </div><!-- /.col -->
</template>

<script>
    export default {
        props: [
            "editing"
        ],
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
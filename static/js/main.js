const GoRedisWebAdmin = {
    data() {
        return {
            enableAutoRefreshDBTree: true,
            dbTree: [],
            activeDBNodes: [],
            invalidTokenDialogVisible: false,
            operateForm: {
                host: null,
                index: null,
                type: null,
                action: null,
                key: null,
                field: null,
                value: null,
                score: null,
                min: null,
                max: null,
                ttl: null
            },
            operateResponse: null,
            operateInputStatus: {
                keyInputVisible: false,
                fieldInputVisible: false,
                valueInputVisible: false,
                scoreInputVisible: false,
                rangeInputVisible: false,
                valueInputPlaceholder: ""
            },
            supportedActions: [],
            userName: null
        }
    },
    methods: {
        getDBTree() {
            let that = this;
            axios({
                method: 'get',
                url: '/rds/db-tree',
            }).then(function (res) {
                if (res.data.code === 0) {
                    let hosts = res.data.data.hosts;
                    that.$data.dbTree = hosts;
                } else if (res.data.code === -6) {
                    that.$data.invalidTokenDialogVisible = true;
                } else {
                    alert(res.data.msg);
                }
            })
        },
        refreshDBTree() {
            let that = this;
            axios({
                method: 'get',
                url: '/rds/refresh-db-tree',
            }).then(function (res) {
                if (res.data.code === 0) {
                    let hosts = res.data.data.hosts;
                    that.$data.dbTree = hosts;
                } else if (res.data.code === -6) {
                    that.$data.invalidTokenDialogVisible = true;
                } else {
                    alert(res.data.msg);
                }
            });
        },
        goToLoginPage: function () {
            window.location.href = '/view/login';
        },
        selectPage: function (node, index, type) {
            let that = this;
            that.$data.operateForm.host = node;
            that.$data.operateForm.index = index;
            that.$data.operateForm.type = type;
        },
        doOperate: function () {
            let that = this;
            let params = that.$data.operateForm;
            if (params.host === null || params.index === null || params.action === null || params.key === null) {
                alert('Not enough parameters');
                return;
            }
            params.index = params.index + "";

            axios({
                method: 'post',
                url: '/rds/operate',
                data: params,
            }).then(function (res) {
                if (res.data.code === 0) {
                    that.$data.operateResponse = res.data.data;
                } else if (res.data.code === -6) {
                    that.$data.invalidTokenDialogVisible = true;
                } else {
                    that.$data.operateResponse = res.data.msg;
                }
            });
        },
        doReset: function () {
            let that = this;
            that.removeAllInputValue();
            that.resetInputComponent();
            that.clearOperationResponse();
        },
        removeAllInputValue: function () {
            let that = this;
            that.$data.operateForm.host = null;
            that.$data.operateForm.index = null;
            that.$data.operateForm.type = null;
            that.$data.operateForm.action = null;
            that.$data.operateForm.key = null;
            that.$data.operateForm.field = null;
            that.$data.operateForm.value = null;
            that.$data.operateForm.ttl = null;
            that.$data.operateForm.score = null;
            that.$data.operateForm.min = null;
            that.$data.operateForm.max = null;
        },
        resetInputComponent: function () {
            let that = this;
            that.$data.operateInputStatus.keyInputVisible = false;
            that.$data.operateInputStatus.fieldInputVisible = false;
            that.$data.operateInputStatus.scoreInputVisible = false;
            that.$data.operateInputStatus.valueInputVisible = false;
            that.$data.operateInputStatus.rangeInputVisible = false;
            that.$data.operateInputStatus.valueInputPlaceholder = "";
        },
        onActionChange: function (val) {
            let that = this;
            let supportedActions = that.$data.supportedActions;
            that.resetInputComponent();

            if (val !== null) {
                for (let i = 0; i < supportedActions.length; i++) {
                    let curAction = supportedActions[i];
                    if (curAction.action === val) {
                        for (let j = 0; j < curAction.requiredParams.length; j++) {
                            let param = curAction.requiredParams[j];
                            if (param === 'Key') {
                                that.$data.operateInputStatus.keyInputVisible = true;
                            }
                            if (param === 'Field') {
                                that.$data.operateInputStatus.fieldInputVisible = true;
                            }
                            if (param === 'Score') {
                                that.$data.operateInputStatus.scoreInputVisible = true;
                            }
                            if (param === 'Min' || param === 'Max') {
                                that.$data.operateInputStatus.rangeInputVisible = true;
                            }
                            if (param === 'Value') {
                                that.$data.operateInputStatus.valueInputVisible = true;
                                that.$data.operateInputStatus.valueInputPlaceholder = curAction.tips;
                            }
                        }
                    }
                }
            }
        },
        getSupportedActions: function () {
            let that = this;
            axios({
                method: 'get',
                url: '/rds/supportedActions',
            }).then(function (res) {
                if (res.data.code === 0) {
                    that.$data.supportedActions = res.data.data.actions;
                } else if (res.data.code === -6) {
                    that.$data.invalidTokenDialogVisible = true;
                } else {
                    alert(res.data.msg);
                }
            });
        },
        clearOperationResponse: function () {
            let that = this;
            that.$data.operateResponse = null;
        },
        doExit: function () {
            axios({
                method: 'post',
                url: '/auth/logout',
            }).then(function (res) {});
            localStorage.removeItem("go-redis-user");
            localStorage.removeItem("go-redis-token");
            window.location.href = "/view/login";
        }
    },
    template: `
    <div style="height: 850px">
        <el-menu class="el-menu-demo" mode="horizontal" background-color="#8895ff" text-color="#fff" style="margin: -8px; height: 45px">
            <span>
                <a href="https://gitee.com/PlusOnline1995">
                    <img title="Go Redis WebAdmin" src="/static/img/logo.png" style="height: 40px">
                </a>
            </span>
            <el-sub-menu index="1">
                <template #title>{{userName}}</template>
                <el-menu-item @click="doExit">Logout</el-menu-item>
            </el-sub-menu>
        </el-menu>
        <el-row :gutter="10" style="margin-top:10px">
            <el-col :span="3"><div></div></el-col>
            <el-col :span="6"><div class="grid-content bg-purple">
                <el-collapse v-model="activeDBNodes" style="padding:10px; overflow: scroll; height:780px" accordion>
                    <el-collapse-item v-for="dbNode in dbTree" :title="dbNode.host" :name="dbNode.host">
                        <ul class="infinite-list" style="overflow: auto">
                            <li v-for="page in dbNode.partitions" class="infinite-list-item" @click="selectPage(dbNode.host, page.index,dbNode.type)">partition：{{page.index}}，size：{{page.DBSize}}</li>
                        </ul>
                    </el-collapse-item>
                </el-collapse>
                <el-button @click="refreshDBTree" size="small" style="margin: 5px;">Refresh</el-button>
                <el-switch v-model="enableAutoRefreshDBTree" size="small" style="margin: 5px;" class="ml-2" active-color="#13ce66" inactive-color="#ff4949" active-text="Auto Refresh"/>
                <br />
            </div></el-col>
            <el-col :span="12" style="height:850px; max-height: 850px"><div class="grid-content bg-purple">
                <el-form label-width="80px" style="padding:10px">
                    <el-row>
                        <el-col :span="8">
                            <el-form-item label="database">
                                <el-input v-model="operateForm.host" readonly />
                            </el-form-item>
                        </el-col>
                        <el-col :span="8">
                            <el-form-item label="partition">
                                <el-input v-model="operateForm.index" readonly />
                            </el-form-item>
                        </el-col>
                        <el-col :span="8">
                            <el-form-item label="type">
                                <el-input v-model="operateForm.type" readonly />
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="8">
                            <el-form-item label="action">
                                <el-select v-model="operateForm.action" placeholder="action" @change="onActionChange">
                                    <el-option v-for="action in supportedActions" :label="action.action" :value="action.action"></el-option>
                                </el-select>
                            </el-form-item>
                        </el-col>
                        <el-col :span="8" v-show="operateInputStatus.keyInputVisible">
                            <el-form-item label="key">
                                <el-input v-model="operateForm.key" type="text"/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="8" v-show="operateInputStatus.fieldInputVisible">
                            <el-form-item label="field">
                                <el-input v-model="operateForm.field" type="text"/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="8" v-show="operateInputStatus.scoreInputVisible">
                            <el-form-item label="score">
                                <el-input v-model="operateForm.score" type="text"/>
                            </el-form-item>
                        </el-col>
                        <el-col :span="4" v-show="operateInputStatus.rangeInputVisible">
                            <el-form-item label="min">
                                <el-input type="text" v-model="operateForm.min"></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="4" v-show="operateInputStatus.rangeInputVisible">
                            <el-form-item label="max">
                                <el-input type="text" v-model="operateForm.max"></el-input>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-form-item label="value" v-show="operateInputStatus.valueInputVisible">
                        <el-input type="text" v-model="operateForm.value" :placeholder="operateInputStatus.valueInputPlaceholder"/>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="warning" @click="doReset">reset</el-button>
                        <el-button type="primary" @click="doOperate">execute</el-button>
                    </el-form-item>
                    <el-form-item label="response">
                        <el-input type="textarea" readonly v-model="operateResponse" />
                    </el-form-item>
                </el-form>
            </div></el-col>
            <el-col :span="3"><div></div></el-col>
        </el-row>
        <el-dialog v-model="invalidTokenDialogVisible" title="Tips" :show-close=false width="30%" :close-on-click-modal="false">
            <span>Invalid token, please login.</span>
            <template #footer>
                <span class="dialog-footer">
                    <el-button @click="goToLoginPage">Ok</el-button>
                </span>
            </template>
        </el-dialog>
    </div>
    `,
    mounted: function () {
        let that = this;

        if (localStorage.getItem("go-redis-token") === undefined) {
            that.$data.invalidTokenDialogVisible = true;
        }

        that.$data.userName = localStorage.getItem("go-redis-user");
        that.getDBTree();
        that.getSupportedActions();
        setInterval(function () {
            if (that.enableAutoRefreshDBTree) {
                that.getDBTree();
            }
        }, 30000);
    }
}
const app = Vue.createApp(GoRedisWebAdmin);
app.use(ElementPlus);
app.mount('#readingo-tpl');

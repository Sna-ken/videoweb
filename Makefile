MODULE := github.com/Sna-ken/videoweb
CMD := cmd
IDL_PATH := idl
KITEX_GEN := kitex_gen
KITEX_TEMPLATE := template/kitex
API_DIR := api
HERTZ_LAYOUT := template/hertz/layout.yaml

# 运行全部单元测试；Mockey 需要关闭编译优化和函数内联
.PHONY: test
test:
	go test -gcflags="all=-N -l" ./...

# 生成基于 Kitex 的业务代码，在新建业务时使用
# 示例：make kitex-gen-auth、make kitex-gen-user
.PHONY: kitex-gen-%
kitex-gen-%:
	kitex \
	-gen-path $(KITEX_GEN) \
	-module "$(MODULE)" \
	-type thrift \
	-thrift template=slim \
	-template-dir $(KITEX_TEMPLATE) \
	$(IDL_PATH)/$*.thrift
	chmod +x rpc/$*/build.sh rpc/$*/script/bootstrap.sh
	go mod tidy

# 更新根目录 kitex_gen 中的对应模块，不影响 cmd 中的业务代码
# 示例：make kitex-update-auth、make kitex-update-user
.PHONY: kitex-update-%
kitex-update-%:
	kitex \
	-gen-path $(KITEX_GEN) \
	-module "$(MODULE)" \
	-thrift template=slim \
	$(IDL_PATH)/$*.thrift

# 首次生成 Hertz API 脚手架
# model、handler、router、script 生成到 api，入口生成到 cmd/api
.PHONY: hertz-gen-api
hertz-gen-api:
	@if [ -f .hz ]; then \
		echo "Hertz 项目已存在，改为执行 hertz-update-api"; \
		$(MAKE) hertz-update-api; \
	else \
		hz new \
		-module "$(MODULE)" \
		-out_dir . \
		-handler_dir $(API_DIR)/handler \
		-model_dir $(API_DIR)/model \
		-router_dir $(API_DIR)/router \
		-idl $(IDL_PATH)/api.thrift \
		-customize_layout $(HERTZ_LAYOUT) \
		-t template=slim && \
		go mod tidy; \
	fi

# IDL 变更后更新 Hertz API 代码
# router_dir 会从首次生成的根目录 .hz 配置中读取
.PHONY: hertz-update-api
hertz-update-api:
	hz update \
	-module "$(MODULE)" \
	-out_dir . \
	-handler_dir $(API_DIR)/handler \
	-model_dir $(API_DIR)/model \
	-idl $(IDL_PATH)/api.thrift \
	-t template=slim

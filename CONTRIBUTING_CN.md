# 贡献指南 - GO-POSTOFFICE

[English](CONTRIBUTING.md) 

我们欢迎您为 GO-POSTOFFICE 项目做出贡献！本文档提供了项目贡献的指南。

## 开始

1. 在 GitHub 上 fork 本仓库。
2. 将您的 fork 克隆到本地：
   ```
   git clone https://github.com/your-username/GO-POSTOFFICE.git
   ```
3. 为您的功能或 bug 修复创建一个新分支：
   ```
   git checkout -b feature/your-feature-name
   ```

## 进行更改

1. 在您的功能分支中进行更改。
2. 根据需要添加或更新测试。
3. 确保您的代码遵循项目的编码标准。
4. 运行测试套件以确保所有测试都通过。

## 提交更改

1. 将您的更改推送到 GitHub 上的 fork：
   ```
   git push origin feature/your-feature-name
   ```
2. 从您的 fork 创建一个 pull request 到主 GO-POSTOFFICE 仓库。
3. 在 pull request 描述中描述您的更改。
4. 在 pull request 描述中链接任何相关的 issues。

## 代码审查

1. 维护者将审查您的 pull request。
2. 根据要求进行任何更改。
3. 根据需要向您的分支推送额外的提交。

## 在您的 Pull Request 被合并后

在您的 pull request 被合并后，您可以安全地删除您的分支并从主（上游）仓库拉取更改：

1. 删除 GitHub 上的远程分支：
   ```
   git push origin --delete feature/your-feature-name
   ```
2. 检出主分支：
   ```
   git checkout main
   ```
3. 删除本地分支：
   ```
   git branch -d feature/your-feature-name
   ```
4. 用最新的上游版本更新您的主分支：
   ```
   git pull upstream main
   ```

## 报告 Bug

1. 通过在 GitHub 的 Issues 中搜索，确保该 bug 尚未被报告。
2. 如果您找不到解决该问题的开放 issue，请开一个新的。
3. 包括一个清晰的标题和详细描述，尽可能多的相关信息，以及一个代码示例或可执行的测试用例，演示未发生的预期行为。

## 建议增强功能

1. 开一个带有清晰标题和详细描述的新 issue。
2. 解释为什么这个增强功能对大多数 GO-POSTOFFICE 用户有用。
3. 提供如何使用该功能的示例。

感谢您为 GO-POSTOFFICE 做出贡献！
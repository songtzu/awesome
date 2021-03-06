

# 开发设计哲学

## mod is evil
我认为go的mod，让go开发娱乐圈化，任何代码，都可以不加思考的廉价引入进来。一个公司的业务拆分给不同的人处理，可能同一个功能函数，会引用两个不同的第三方包的实现。

认真对待自己的代码工程，这是一个严肃的原则问题。

更有甚者，三流团队开发者，能够整出如下局面：
```
a import b mod
a import c mod
b import h mod v1
c import h mod v2
h mod v1与v2不兼容
```
开发者对他们引入的第三方包，可以被轻松廉价的删除这个事实，毫无认知。

mod提供了一个极便利的手段，让那些不自律、不愿意思考的人更加堕落。犹如把核武器开关交付给心智不健全的人。

引入mod犹如繁华的交叉路口移除了交通信号灯，谁野蛮谁说了算。保留一点偷懒门槛，让那些不合格的开发者能够进而远之，让合格的开发者想偷懒的时候能够稍加思考。
## pure and simple.
我们的世界太复杂了。与别人的交流协作也太复杂。要做到真正的keep it simple and easy to use真的很难。

明确的规则和强力的约束，搭配灵活的使用，是我追求的纯粹与简单。

### 冗余的妥协
在anet包，可能会用到某些切片处理的工具函数。可能其他第三方库有提供，甚至awesome仓库自己都有提供，我还是会考虑冗余一份代码在anet包里面。

这样，其他人复制anet文件夹，就可以完整的使用这个文件夹的代码。而不会遇到引用缺失的情况。

我践行严格执行的纯粹派做法。当我可以放开一个口子，在一个很纯粹的网络模块，引入本仓库的其他包，甚至是其他第三方仓库的代码，一切都会迈向注定的不停妥协的结局。

* 之所以放在此package中，是考虑到mq包可能会独立被引用或者移植的情形，如果引用了本仓库中的其他包里面的内容，则无法做到简单的复制文件夹即可用。

## 好的代码必然容易写单元测试代码
我经常见到一些代码，耦合各种模块。例如，数据库包，直接引入配置包，毕竟数据库的地址在配置文件里面，而配置的初始化，往往在main函数中调用的。
我们要给数据库包做单元测试，还需要先初始化配置包。

实际上，数据库包明明可以在初始化的函数中，传入相关的配置。这样在写单元测试的时候，可以不用关系配置文件包的代码和初始化。

上面说的还是极其简化的局面。我们把0.9分的代码重复写100次，整个系统的得分就趋于0分了。


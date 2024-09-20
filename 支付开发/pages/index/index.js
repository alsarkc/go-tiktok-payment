const app = getApp()

Page({
  data: {

  },
  onLoad: function () {
    console.log('Welcome to Mini Code')
  },


  pay() {

    tt.request({
      url: "http://127.0.0.1:8080/preorder",
      success: (res) => {
        console.log(res.data)
        let order_id = res.data.data.order_id
        let order_token = res.data.data.order_token



        tt.pay({
          orderInfo: {
            order_id: order_id,
            order_token: order_token,
          },
          service: 5,
          success(res1) {
            if (res1.code == 0) {
              // 支付成功处理逻辑，只有res.code=0时，才表示支付成功
              // 但是最终状态要以商户后端结果为准

              tt.showToast({
                title: '支付成功',
              });

            }
          },
          fail(res1) {
            // 调起收银台失败处理逻辑
          },
        });






      }
    });



  }



})
package middleware

//
//import (
//	"context"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"time"
//)
//func (w *timeoutWriter) markTimeout() {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//
//	w.timedOut = true
//	w.buf.Reset()
//}
//
//func Timeout(d time.Duration) gin.HandlerFunc {
//	if d <= 0 {
//		d = 5 * time.Second
//	}
//
//	return func(c *gin.Context) {
//		// 1️⃣ 创建可取消 context
//		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
//		defer cancel()
//		c.Request = c.Request.WithContext(ctx)
//
//		// 2️⃣ 替换 ResponseWriter
//		tw := &timeoutWriter{
//			ResponseWriter: c.Writer,
//		}
//		c.Writer = tw
//
//		finished := make(chan struct{})
//		panicChan := make(chan any, 1)
//
//		// 3️⃣ 业务放 goroutine 执行
//		go func() {
//			defer func() {
//				if p := recover(); p != nil {
//					panicChan <- p
//				}
//			}()
//			c.Next()
//			close(finished)
//		}()
//
//		// 4️⃣ 三选一：完成 / panic / timeout
//		select {
//
//		case <-finished:
//			// ✅ 正常完成：一次性 flush
//			tw.flush()
//			return
//
//		case p := <-panicChan:
//			// ✅ panic 继续往上抛（交给 recovery middleware）
//			panic(p)
//
//		case <-ctx.Done():
//			// ⏰ 超时
//			tw.markTimeout()
//
//			// ⚠️ 用“真实 writer”返回 timeout
//			orig := tw.ResponseWriter
//			if !orig.Written() {
//				c.Writer = orig
//				c.AbortWithStatusJSON(
//					http.StatusGatewayTimeout,
//					Response{
//						Code:    504,
//						Message: "request_timeout",
//					},
//				)
//			}
//			return
//		}
//	}
//}
//
////func Timeout(d time.Duration) gin.HandlerFunc {
////	if d <= 0 {
////		d = 5 * time.Second
////	}
////	return func(c *gin.Context) {
////		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
////		defer cancel()
////		c.Request = c.Request.WithContext(ctx)
////
////		finished := make(chan struct{})
////		panicChan := make(chan any, 1)
////
////		go func() {
////			defer func() {
////				if p := recover(); p != nil {
////					panicChan <- p
////				}
////			}()
////			c.Next()
////			close(finished)
////		}()
////
////		select {
////		case <-finished:
////			return
////		case p := <-panicChan:
////			panic(p)
////		case <-ctx.Done():
////			if !c.Writer.Written() {
////				resp.Fail(c, http.StatusGatewayTimeout, "request_timeout")
////				c.Abort()
////			}
////		}
////	}
////}

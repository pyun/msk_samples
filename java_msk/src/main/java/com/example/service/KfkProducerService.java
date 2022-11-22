package com.example.service;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Component;
import org.springframework.util.concurrent.ListenableFuture;
import org.springframework.util.concurrent.ListenableFutureCallback;

@Component
@Slf4j
public class KfkProducerService {
    @Autowired
    private KafkaTemplate<String, Object> kafkaTemplate;

    @Value("${spring.kafka.bootstrap-servers}")
    private String brokers;

    public void startProducer() {
        log.info("[connect to cluster:{}", brokers);
        long id = 10000000000L;
        int i = 0;
        String message = "test msage";
        while (i < 5) {
            id += 1;
            i += 1;
            send("test-topic1", "[" + id + "]" + message);

            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }

        }
    }

    public void send(String topic, String str) {
        //log.info("[准备发送消息][{}]{}", topic, str.substring(0, 13));
        //发送消息
        try {
            ListenableFuture<SendResult<String, Object>> future = kafkaTemplate.send(topic, str);
            future.addCallback(new ListenableFutureCallback<SendResult<String, Object>>() {

                public void onFailure(Throwable throwable) {
                    //发送失败的处理
                    //log.info(TOPIC_TEST + " - 生产者 发送消息失败：" + throwable.getMessage());
                    log.info("[发送消息失败][{}]{}", topic, throwable.getMessage());
                    //todo ...
                }


                public void onSuccess(SendResult<String, Object> stringObjectSendResult) {
                    //成功的处理
                    log.info("[发送消息成功][{}]{}", topic, stringObjectSendResult.getProducerRecord().value().toString().substring(0, 13));
                    //todo ...
                }
            });
        } catch (Exception e) {
            e.printStackTrace();
        }


    }
}

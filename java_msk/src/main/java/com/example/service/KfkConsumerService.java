package com.example.service;

import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Component;

import java.util.Optional;

/**
 * 消费者
 */
@Component
@Slf4j
public class KfkConsumerService {
    //消费组1
    public static final String TOPIC_GROUP1 = "topic.group1";
    //消费组2
    public static final String TOPIC_GROUP2 = "topic.group2";

    /**
     * 消费者1
     *
     * @param record
     * @param topic
     */
    @KafkaListener(topics = "test-topic1", groupId = TOPIC_GROUP1)
    public void topic_test1(ConsumerRecord<?, ?> record, @Header(KafkaHeaders.RECEIVED_TOPIC) String topic) {

        Optional message = Optional.ofNullable(record.value());
        if (message.isPresent()) {
            Object msg = message.get();
            log.info("[消费者1 消费了]{}", msg.toString().substring(0, 13));
            //ack.acknowledge();
        }
    }


    /**
     * 消费者2
     *
     * @param record
     * @param topic
     */
    //@KafkaListener(topics = "${spring.kafka.topic}", groupId = TOPIC_GROUP2)
    public void topic_test4(ConsumerRecord<?, ?> record, @Header(KafkaHeaders.RECEIVED_TOPIC) String topic) {

        Optional message = Optional.ofNullable(record.value());
        if (message.isPresent()) {
            Object msg = message.get();
            log.info("消费者2 消费了： Topic:" + topic + ",Message:" + msg);
            //ack.acknowledge();
        }
    }
}

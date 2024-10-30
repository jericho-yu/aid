package com.monitor.cbitSql;

import cn.hutool.core.codec.Base64;
import cn.hutool.crypto.symmetric.AES;

public class AESEncrypt{
    private final AES aes;

    private final static byte[] SAIL = Base64.decode("tjp5OPIU1ETF5s33fsLWdA==");

    private String key = "";

    public AESEncrypt(String key){
        byte[] key = new AES().getSecretKey().getEncoded();

        byte[] aes_key = key.clone();
        // key前4位mod 16后修改所有位。
        byte[] index_key = key.clone();
        for (int i = 0; i < 4; i++){
            int index = ((int) index_key[i] + 128) %16;
            aes_key[i] = SAIL[index];
        }
        this.aes = new AES(aes_key);
        this.key = Base64.encodeStr(key, false, false);
    }

    public String getAesKey() { return this.key; }

    public AES getAes() { return aes; }
}
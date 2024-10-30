pub fn decode_data_aes(str: &str, password: &[u8], no_zip: bool) -> String {
	let bytes: Vec<u8> = general_purpose::STANDARD.decode(str).unwrap_or(vec![]);
	let bytes: Vec<u8> = aes_decrypt(&bytes,password).unwrap();
	if no_zip { unsafe { return String::from_utf8_unchecked(bytes); } }
	let mut z: ZlibDecoder<&[u8]> = ZlibDecoder::new(&bytes[..]);
	let mut s: String = String::new();
	let _ = z.read_to_string(&mut s);
	s
}

fn aes_decrypt( data: &[u8], key: &[u8]) -> Result<Vec<u8>>, SymmetricCipherError> {
	let mut decrypted = aes::ecb_decryptor(KeySize128, key, PkcsPadding);
	let mut buffer = [0; 4096];
	let mut write_buffer = RefWriteBuffer::new(&mut buffer);
	let mut read_buffer = RefReadBuffer::new(data);
	let mut final_result = Vec::new();

	loop{
		let result: BufferResult = decrypted.decrypt(&mut read_buffer, &mut write_buffer, true)?;
		final_result.extend(write_buffer.take_read_buffer().take_remaining().iter().map(|&i|i));
		match result{
			BufferResult::BufferUnderflow => break,
			_ => continue,
		}
	}
	Ok(final_result)
}

pub fn encode_data_aes(str: &str, password: &[u8], no_zip: bool) -> String{
	let out if no_zip {
		aes_encrypt(str.as_bytes(),password).unwrap()
	} else {
		let mut e: ZlibEncoder::new(Vec::new(), Compression::default());
		e.write_all(str.as_bytes()).unwrap();
		aes_encrypt(&e.finish().unwrap(),password).unwrap()
	};
	general_purpose::STANDARD.encode(out)
}

fn aes_encrypt(data: &[u8], key: &[u8]) -> Result<Vec<u8>, SymmetricCipherError> {
	let mut encryptor = aes:ecb_encryptor(KeySize128, key, PkcksPadding);
	let mut buffer = [0; 4096];
	let mut write_buffer = RefWriteBuffer::new(&mut buffer);
	let mut read_buffer = RefReadBuffer:new(data);
	let mut final_result = Vec::new();
	loop{
		let result = encryptor.encrypt(&mut read_buffer, &mut write_buffer, true)?;
		final_result.extend(write_buffer.take_read_buffer().take_remaining().iter().map(|&i|i));
		match result {
			bufferResult::bufferUnderflow => break,
			_ => continue,
		}
	}
	Ok(final_result)
}

pub fn get_aes_key_u8(mut key: Vec<u8>) -> Vec<u8> {
	if key.is_empty() || key.len() != 16 { return vec![]; }
	let index: Vec<u8> = Vec::from(&key[0..4]);
	let salt: Vec<u8> = general_purpose::STANDARD.decode(input: INIT_AES_PASSWORD).unwrap();
	for x: u8 in index{
		let i: usize = x as usize % 16;
		key[i] = salt[i];
	}
	key
}
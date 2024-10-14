use actix_web::{web, App, HttpServer, Responder, HttpResponse};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use actix_web::Either;

#[derive(Serialize, Deserialize, Debug, Clone)]
struct Item {
    caption: String,
    weight: f32,
    number: i32,
}

struct ItemStore {
    items: HashMap<String, Item>,
}

impl ItemStore {
    fn new() -> Self {
        ItemStore {
            items: HashMap::new(),
        }
    }

    fn add_item(&mut self, item: Item) {
        self.items.insert(item.caption.clone(), item);
    }

    fn get_item(&self, caption: &str) -> Option<&Item> {
        self.items.get(caption)
    }
}

async fn add_item(item_store: web::Data<Arc<Mutex<ItemStore>>>, item: web::Json<Item>) -> impl Responder {
    let mut item_store_lock = item_store.lock().unwrap();
    println!("{:?}", item);
    item_store_lock.add_item(item.into_inner());
    "Элемент добавлен!"
}

async fn get_item(item_store: web::Data<Arc<Mutex<ItemStore>>>, caption: web::Path<String>) -> impl Responder {
    let item_store_lock = item_store.lock().unwrap();
    if let Some(item) = item_store_lock.get_item(&caption) {
        let item_clone = item.clone(); 
        Either::Left(web::Json(item_clone))
    } else {
        Either::Right(HttpResponse::NotFound().finish())
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let item_store = web::Data::new(Arc::new(Mutex::new(ItemStore::new())));
    HttpServer::new(move || {
        App::new()
            .app_data(item_store.clone())
            .service(web::resource("/item").route(web::post().to(add_item)))
            .service(web::resource("/item/{caption}").route(web::get().to(get_item)))
    })
        .bind("127.0.0.1:8080")?
        .run()
        .await
}




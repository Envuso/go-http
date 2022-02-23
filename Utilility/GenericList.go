package Utilility

//
//
// type GenericList[K comparable|string, V any] struct {
// 	data map[K]V
// }
//
// func (p *GenericList[K, V]) Data() map[K]V {
// 	return p.data
// }
//
// func (p *GenericList[K, V]) SetData(data any) {
// 	p.data = data.(map[K]V)
// }
//
// func (p *GenericList[K, V]) Set(key K, value V) {
// 	p.data[key] = value
// }
//
// func (p *GenericList[K, V]) Has(key K) bool {
// 	_, ok := p.data[key]
// 	return ok
// }
//
// func (p *GenericList[K, V]) Get(key K, defaultVal ...V) V {
// 	val, ok := p.data[key]
//
// 	if ok {
// 		return val
// 	}
//
// 	if len(defaultVal) >= 1 {
// 		return defaultVal[0]
// 	}
//
// 	return nil
// }
//
// func (p *GenericList[K, V]) Keys() []K {
// 	return maps.Keys(p.data)
// }
//
// func (p *GenericList[K, V]) Values() []V {
// 	return maps.Values(p.data)
// }
//
//
//
// func (p *GenericList[K, V]) Clear() {
// 	maps.Clear(p.data)
// }
